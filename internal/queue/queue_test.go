package queue

import (
	"context"
	"testing"
	"time"

	"queue-svc/internal/model"
)

// Проверяем, что задача успешно кладется в очередь и читается обратно
func TestEnqueueAndRead(t *testing.T) {
	q := New(2) // очередь на 2 задачи
	ctx := context.Background()

	task := &model.Task{ID: "1"} // создаем задачу
	if err := q.Enqueue(ctx, task); err != nil {
		t.Fatalf("enqueue failed: %v", err)
	}

	select {
	case got := <-q.Chan():
		if got.ID != "1" {
			t.Errorf("expected id=1, got %s", got.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting task from channel")
	}
}

// Проверяем, что после Close() новые задачи не принимаются
func TestClosePreventsNewEnqueue(t *testing.T) {
	q := New(1)
	q.Close()

	err := q.Enqueue(context.Background(), &model.Task{ID: "2"})
	if err != ErrQueueClosed {
		t.Errorf("expected ErrQueueClosed, got %v", err)
	}
}

// Проверяем, что если контекст отменен - Enqueue не блокируется и возвращает ошибку
func TestEnqueueRespectsContext(t *testing.T) {
	q := New(0)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // сразу отменяем

	if err := q.Enqueue(ctx, &model.Task{ID: "x"}); err == nil {
		t.Fatal("expected context error, got nil")
	}
}
