//go:build go1.21
// +build go1.21

package log

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

// NewSLogChiMiddleware is used to log http request information. It takes
// a pointer to an slog.Logger to use. If `l` is nil, it uses the
// default logger
func NewSLogChiMiddleware(l *slog.Logger) func(http.Handler) http.Handler {
	if l == nil {
		l = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func(start time.Time) {
				l.LogAttrs(
					r.Context(),
					slog.LevelInfo,
					"HTTP Request Served",
					slog.String("proto", r.Proto),
					slog.String("path", r.URL.Path),
					slog.Duration("duration", time.Since(start)),
					slog.Int("status", ww.Status()),
					slog.Int("size", ww.BytesWritten()),
					slog.String("ip", r.RemoteAddr),
				)
			}(time.Now())

			next.ServeHTTP(ww, r)
		})
	}
}
