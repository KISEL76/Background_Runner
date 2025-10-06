package worker

import (
	"context"
	"math/rand"
	"queue-svc/internal/model"
	"queue-svc/internal/queue"
	"queue-svc/internal/store"
	"sync"
	"time"
)

const (
	chance_of_fail = 20
)

type Pool struct {
	n     int
	q     *queue.Queue
	s     *store.Store
	wg    sync.WaitGroup
	base  time.Duration
	cap   time.Duration
	failP int
}

func New(n int, q *queue.Queue, s *store.Store) *Pool {
	return &Pool{
		n:     n,
		q:     q,
		s:     s,
		base:  100 * time.Millisecond,
		cap:   2 * time.Second,
		failP: chance_of_fail,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.n; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			rnd := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))
			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-p.q.Chan():
					if !ok {
						return
					}
					p.process(ctx, rnd, t)
				}
			}
		}(i)
	}
}

func (p *Pool) process(ctx context.Context, rnd *rand.Rand, t *model.Task) {
	p.s.UpdateStatus(t.ID, model.StatusRunning, "")

	for {
		work := time.Duration(100+rnd.Intn(401)) * time.Millisecond
		select {
		case <-ctx.Done():
			return
		case <-time.After(work):
		}

		if rnd.Intn(100) >= p.failP {
			p.s.UpdateStatus(t.ID, model.StatusDone, "")
			return
		}

		attempt := p.s.IncAttempts(t.ID)
		if attempt > t.MaxRetries {
			p.s.UpdateStatus(t.ID, model.StatusFailed, "exceeded max retries")
			return
		}
		sleep := p.backoffWithJitter(rnd, t.Attempts)

		select {
		case <-ctx.Done():
			return
		case <-time.After(sleep):
		}
	}
}

func (p *Pool) backoffWithJitter(rnd *rand.Rand, attempt int) time.Duration {
	max := p.base * (1 << (attempt - 1))
	if max > p.cap {
		max = p.cap
	}
	if max <= 0 {
		max = p.base
	}
	return time.Duration(rnd.Int63n(int64(max)))
}

func (p *Pool) Wait() { p.wg.Wait() }
