package handlers

import "net/http"

func NewMux(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /enqueue", h.Enqueue)
	return mux
}
