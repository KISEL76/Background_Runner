package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"queue-svc/internal/handlers"
	"queue-svc/internal/model"
	"queue-svc/internal/queue"
	"queue-svc/internal/store"
	"queue-svc/internal/worker"
)

const N = 10

func Test_FullSystemFlow(t *testing.T) {
	// Сетап всего необходимого
	st := store.New()
	q := queue.New(32)
	pool := worker.New(4, q, st)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool.Start(ctx)

	h := &handlers.Handler{Q: q, S: st}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/enqueue", h.Enqueue)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	// healthcheck
	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("healthz request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Загоняем несколько задач в очередь
	for i := 0; i < N; i++ {
		task := map[string]interface{}{
			"id":          "task-" + itoa(i),
			"payload":     "data",
			"max_retries": 2,
		}
		body, _ := json.Marshal(task)
		r, err := http.Post(srv.URL+"/enqueue", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("enqueue failed: %v", err)
		}
		if r.StatusCode != http.StatusAccepted {
			t.Fatalf("expected 202, got %d", r.StatusCode)
		}
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
	}

	// Ждем выполнения задач
	waitUntil(t, 2*time.Second, func() bool {
		snap := st.Snapshot()
		done, failed := 0, 0
		for _, v := range snap {
			if v.Status == model.StatusDone {
				done++
			}
			if v.Status == model.StatusFailed {
				failed++
			}
		}
		return done+failed == N
	})

	// Финальная проверка, что все задачи завершились
	snap := st.Snapshot()
	done, failed := 0, 0
	for _, v := range snap {
		switch v.Status {
		case model.StatusDone:
			done++
		case model.StatusFailed:
			failed++
		default:
			t.Errorf("task %s stuck in %v", v.ID, v.Status)
		}
	}
	t.Logf("Integration test result: done=%d failed=%d", done, failed)
	if done+failed != N {
		t.Fatalf("not all tasks finished")
	}
}

func waitUntil(t *testing.T, timeout time.Duration, ok func() bool) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if ok() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timeout after %v", timeout)
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n := i; n > 0; n /= 10 {
		pos--
		buf[pos] = byte('0' + n%10)
	}
	return string(buf[pos:])
}
