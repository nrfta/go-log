package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
)

// Fields Type to pass when we want to call WithFields for structured logging
// it matches the old logrus Fields type for convenience.
type Fields map[string]interface{}

var (
	logger *slog.Logger
)

func init() {
	New(false, "info")
}

// New creates a new slog.Logger with the given configuration. If a writer is
// provided it will be used as the output destination.
func New(isJSONFormatted bool, logLevel string, out ...io.Writer) *slog.Logger {
	w := io.Writer(os.Stdout)
	if len(out) > 0 && out[0] != nil {
		w = out[0]
	}

	opts := &slog.HandlerOptions{Level: getLevel(logLevel)}
	var h slog.Handler
	if isJSONFormatted {
		h = slog.NewJSONHandler(w, opts)
	} else {
		h = slog.NewTextHandler(w, opts)
	}

	log := slog.New(h)
	logger = log
	return log
}

// GetLogger returns the package logger.
func GetLogger() *slog.Logger { return logger }

// RequestLogger creates a logger with the request ID on it.
func RequestLogger(ctx context.Context) *slog.Logger {
	return logger.With(slog.String("requestID", middleware.GetReqID(ctx)))
}

// Helper to convert string log levels into slog.Level values.
func getLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error", "fatal", "panic":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// helper to convert Fields to slog attributes
func toAttrs(fields Fields) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	return attrs
}

func Info(args ...interface{})              { logger.Info(fmt.Sprint(args...)) }
func Infof(msg string, args ...interface{}) { logger.Info(fmt.Sprintf(msg, args...)) }
func InfoWithFields(fields Fields, args ...interface{}) {
	logger.With(toAttrs(fields)...).Info(fmt.Sprint(args...))
}
func InfoWithFieldsf(fields Fields, msg string, args ...interface{}) {
	logger.With(toAttrs(fields)...).Info(fmt.Sprintf(msg, args...))
}

func Debug(args ...interface{})              { logger.Debug(fmt.Sprint(args...)) }
func Debugf(msg string, args ...interface{}) { logger.Debug(fmt.Sprintf(msg, args...)) }
func DebugWithFields(fields Fields, args ...interface{}) {
	logger.With(toAttrs(fields)...).Debug(fmt.Sprint(args...))
}
func DebugWithFieldsf(fields Fields, msg string, args ...interface{}) {
	logger.With(toAttrs(fields)...).Debug(fmt.Sprintf(msg, args...))
}

func Error(args ...interface{})              { logger.Error(fmt.Sprint(args...)) }
func Errorf(msg string, args ...interface{}) { logger.Error(fmt.Sprintf(msg, args...)) }
func ErrorWithFields(fields Fields, args ...interface{}) {
	logger.With(toAttrs(fields)...).Error(fmt.Sprint(args...))
}
func ErrorWithFieldsf(fields Fields, msg string, args ...interface{}) {
	logger.With(toAttrs(fields)...).Error(fmt.Sprintf(msg, args...))
}

func NewError(args ...interface{}) error {
	Error(args...)
	return fmt.Errorf("%s", fmt.Sprint(args...))
}
func NewErrorf(msg string, args ...interface{}) error {
	Errorf(msg, args...)
	return fmt.Errorf(msg, args...)
}
func NewErrorWithFields(fields Fields, args ...interface{}) error {
	ErrorWithFields(fields, args...)
	return fmt.Errorf("%s", fmt.Sprint(args...))
}
func NewErrorWithFieldsf(fields Fields, msg string, args ...interface{}) error {
	ErrorWithFieldsf(fields, msg, args...)
	return fmt.Errorf(msg, args...)
}

func Warn(args ...interface{})              { logger.Warn(fmt.Sprint(args...)) }
func Warnf(msg string, args ...interface{}) { logger.Warn(fmt.Sprintf(msg, args...)) }
func WarnWithFields(fields Fields, args ...interface{}) {
	logger.With(toAttrs(fields)...).Warn(fmt.Sprint(args...))
}
func WarnWithFieldsf(fields Fields, msg string, args ...interface{}) {
	logger.With(toAttrs(fields)...).Warn(fmt.Sprintf(msg, args...))
}

func Fatal(args ...interface{}) {
	logger.Error(fmt.Sprint(args...))
	os.Exit(1)
}
func Fatalf(msg string, args ...interface{}) {
	logger.Error(fmt.Sprintf(msg, args...))
	os.Exit(1)
}
func FatalWithFields(fields Fields, args ...interface{}) {
	logger.With(toAttrs(fields)...).Error(fmt.Sprint(args...))
	os.Exit(1)
}
func FatalWithFieldsf(fields Fields, msg string, args ...interface{}) {
	logger.With(toAttrs(fields)...).Error(fmt.Sprintf(msg, args...))
	os.Exit(1)
}

func Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	logger.Error(msg)
	panic(msg)
}
func Panicf(msg string, args ...interface{}) {
	formatted := fmt.Sprintf(msg, args...)
	logger.Error(formatted)
	panic(formatted)
}
func PanicWithFields(fields Fields, args ...interface{}) {
	msg := fmt.Sprint(args...)
	logger.With(toAttrs(fields)...).Error(msg)
	panic(msg)
}
func PanicWithFieldsf(fields Fields, msg string, args ...interface{}) {
	formatted := fmt.Sprintf(msg, args...)
	logger.With(toAttrs(fields)...).Error(formatted)
	panic(formatted)
}

// ServerLogger logs the start and end of each request with useful metadata.
func ServerLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				logger.With(
					slog.String("proto", r.Proto),
					slog.String("path", r.URL.Path),
					slog.Duration("duration", time.Since(t1)),
					slog.Int("status", ww.Status()),
					slog.Int("size", ww.BytesWritten()),
					slog.String("ip", r.RemoteAddr),
					slog.String("requestID", middleware.GetReqID(r.Context())),
				).Info("Request Served")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// Writer exposes a writer from the logger's handler when possible. If the
// underlying handler does not support retrieving a writer, nil is returned.
func Writer() *io.PipeWriter            { return nil }
func WriterLevel(string) *io.PipeWriter { return nil }
