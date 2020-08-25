package log

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

//Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// Logger interface
type Logger interface {
	Info(args ...interface{})
	Infof(message string, args ...interface{})
	Debug(args ...interface{})
	Debugf(message string, args ...interface{})
	Error(args ...interface{})
	Errorf(message string, args ...interface{})
	Warn(args ...interface{})
	Warnf(message string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(message string, args ...interface{})
	Panic(args ...interface{})
	Panicf(message string, args ...interface{})
	Writer() *io.PipeWriter
}

var (
	logger *logrus.Logger
)

func init() {
	New(false, "info")
}

// New - Creates a new instance of logrus with customized configuration
func New(isProduction bool, logLevel string) *logrus.Logger {
	var formatter logrus.Formatter

	formatter = &logrus.TextFormatter{
		ForceColors:            true,
		DisableLevelTruncation: true,
	}

	if isProduction {
		formatter = &logrus.JSONFormatter{}
	}
	log := logrus.New()
	log.SetFormatter(formatter)

	switch logLevel {
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	}

	logger = log
	return log
}

// RequestLogger creates a logger with the request ID on it
func RequestLogger(ctx context.Context) Logger {
	return logger.WithFields(logrus.Fields{
		"requestID": middleware.GetReqID(ctx),
	})
}

func Writer() *io.PipeWriter {
	return logger.Writer()
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(message string, args ...interface{}) {
	logger.Infof(message, args...)
}

func InfoWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Info(args...)
}

func InfoWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Infof(message, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(message string, args ...interface{}) {
	logger.Debugf(message, args...)
}

func DebugWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Debug(args...)
}

func DebugWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Debugf(message, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(message string, args ...interface{}) {
	logger.Errorf(message, args...)
}

func ErrorWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Error(args...)
}

func ErrorWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Errorf(message, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(message string, args ...interface{}) {
	logger.Warnf(message, args...)
}

func WarnWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Warn(args...)
}

func WarnWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Warnf(message, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(message string, args ...interface{}) {
	logger.Fatalf(message, args...)
}

func FatalWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Fatal(args...)
}

func FatalWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Fatalf(message, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(message string, args ...interface{}) {
	logger.Panicf(message, args...)
}

func PanicWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Panic(args...)
}

func PanicWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Panicf(message, args...)
}

// ServerLogger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func ServerLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				logger.WithFields(logrus.Fields{
					"proto":     r.Proto,
					"path":      r.URL.Path,
					"duration":  time.Since(t1),
					"status":    ww.Status(),
					"size":      ww.BytesWritten(),
					"ip":        r.RemoteAddr,
					"requestID": middleware.GetReqID(r.Context()),
				}).Info("Request Served")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
