package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"queue-svc/internal/config"
	"queue-svc/internal/handlers"
	"queue-svc/internal/queue"
	"queue-svc/internal/shutdown"
	"queue-svc/internal/store"
	"queue-svc/internal/worker"
	"syscall"
	"time"
)

const (
	address = ":8080"
)

func main() {
	cfg := config.Load()
	s := store.New()
	q := queue.New(cfg.QueueSize)
	h := handlers.NewHandler(q, s)
	mux := handlers.NewMux(h)

	server := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("listening on %s (workers=%d queue=%d)", address, cfg.Workers, cfg.QueueSize)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	p := worker.New(cfg.Workers, q, s)
	p.Start(ctx)

	<-ctx.Done()
	log.Println("shutting down: stop accept...")
	shutdown.GracefulHTTP(server, 10*time.Second)
	q.Close()
	p.Wait()
	log.Println("stopped")
}
