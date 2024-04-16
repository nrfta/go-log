package log

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
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

type VariablesScrubber interface {
	Scrub(map[string]any) map[string]any
}

type noopVariablesScrubber struct{}

var _ VariablesScrubber = (*noopVariablesScrubber)(nil)

func (noopVariablesScrubber) Scrub(vars map[string]any) map[string]any {
	return nil
}

// NewSLogGraphQLResponseMiddleware is used to log GraphQL requests and responses.
func NewSLogGraphQLResponseMiddleware(l *slog.Logger, s VariablesScrubber) graphql.ResponseMiddleware {
	if l == nil {
		l = slog.Default()
	}

	if s == nil {
		s = noopVariablesScrubber{}
	}

	return func(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
		var (
			start = time.Now()
			res   = next(ctx)
			oc    = graphql.GetOperationContext(ctx)
		)

		if !strings.Contains(oc.RawQuery, "__ApolloGetServiceDefinition__") {
			l.LogAttrs(
				ctx,
				slog.LevelInfo,
				"GraphQL Request Served",
				slog.Group(
					"graphql",
					slog.Group(
						"req",
						slog.String("query", oc.RawQuery),
						slog.Any("variables", s.Scrub(oc.Variables)),
					),
					slog.Group(
						"res",
						slog.String("errors", res.Errors.Error()),
					),
					slog.Duration("duration", time.Since(start)),
				),
			)
		}

		return res
	}
}
