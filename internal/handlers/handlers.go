package handlers

import (
	"encoding/json"
	"net/http"
	"queue-svc/internal/model"
	"queue-svc/internal/queue"
	"queue-svc/internal/store"
)

type Handler struct {
	Q *queue.Queue
	S *store.Store
}

func NewHandler(q *queue.Queue, s *store.Store) *Handler {
	return &Handler{
		Q: q,
		S: s,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Enqueue(w http.ResponseWriter, r *http.Request) {
	var t model.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if t.ID == "" {
		http.Error(w, "there is no id", http.StatusBadRequest)
		return

	}
	if t.MaxRetries < 0 {
		t.MaxRetries = 0
	}
	t.Status = model.StatusQueued
	h.S.Put(&t)

	if err := h.Q.Enqueue(r.Context(), &t); err != nil {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "queued", "id": t.ID})
}
