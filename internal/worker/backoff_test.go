package worker

import (
	"math/rand"
	"testing"
	"time"
)

// Проверяем, что backoffWithJitter выдает значения в диапазоне и они разные,
// и что среднее значение увеличивается при росте количества попыток attempt.
func TestBackoffWithJitter_RangeAndVariety(t *testing.T) {
	p := New(1, nil, nil)

	// Для скорости теста уменьшаем базовое и максимальное время задержки ретрая
	p.base = 10 * time.Millisecond
	p.cap = 40 * time.Millisecond

	var prevAvg time.Duration

	// Проверяем несколько попыток подряд
	for attempt := 1; attempt <= 3; attempt++ {
		var values []time.Duration
		var sum time.Duration

		// Верхняя граница задержки
		max := p.base * (1 << (attempt - 1))
		if max > p.cap {
			max = p.cap
		}

		// Один генератор для всей серии
		rnd := rand.New(rand.NewSource(time.Now().UnixNano() + int64(attempt)))

		// Берем 20 выборок
		for i := 0; i < 20; i++ {
			d := p.backoffWithJitter(rnd, attempt)
			// Проверка диапазона
			if d < 0 || d > max {
				t.Fatalf("attempt=%d backoff out of range: got=%v max=%v", attempt, d, max)
			}
			values = append(values, d)
			sum += d
		}

		avg := sum / time.Duration(len(values))
		t.Logf("attempt=%d delays: %v (avg=%v)", attempt, values, avg)

		// Проверяем, что не все значения одинаковые
		allSame := true
		first := values[0]
		for _, v := range values[1:] {
			if v != first {
				allSame = false
				break
			}
		}
		if allSame {
			t.Fatalf("jitter seems constant on attempt=%d: %v", attempt, values)
		}

		// Проверяем, что среднее значение растёт (экспонента работает)
		if attempt > 1 && avg <= prevAvg {
			t.Fatalf("expected avg to grow: prev=%v, current=%v", prevAvg, avg)
		}
		prevAvg = avg
	}
}
