package shutdown

import (
	"context"
	"net/http"
	"time"
)

func GracefulHTTP(srv *http.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
