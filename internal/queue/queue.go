package queue

import (
	"context"
	"errors"
	"queue-svc/internal/model"
	"sync"
)

var ErrQueueClosed = errors.New("queue closed")

type Queue struct {
	ch     chan *model.Task
	mu     sync.RWMutex
	closed bool
}

func New(queueSize int) *Queue {
	return &Queue{
		ch:     make(chan *model.Task, queueSize),
		closed: false,
	}
}

func (q *Queue) Enqueue(ctx context.Context, t *model.Task) error {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.closed {
		return ErrQueueClosed
	}

	select {
	case q.ch <- t:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.closed {
		q.closed = true
		close(q.ch)
	}
}

func (q *Queue) Chan() <-chan *model.Task {
	return q.ch
}
