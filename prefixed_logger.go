package log

import (
	"fmt"
	"log/slog"

	"github.com/neighborly/go-errors"
)

// PrefixedLogger wraps a slog.Logger and prefixes all messages with a given
// string.
type PrefixedLogger struct {
	Prefix         string
	LoggerInstance *slog.Logger
}

// NewPrefixedLogger returns a PrefixedLogger using the provided slog.Logger. If
// nil is given, a new JSON formatted info-level logger is created.
func NewPrefixedLogger(prefix string, instance *slog.Logger) PrefixedLogger {
	var l *slog.Logger
	if instance == nil {
		l = New(true, "info")
	} else {
		l = instance
	}
	return PrefixedLogger{Prefix: prefix, LoggerInstance: l}
}

func (l *PrefixedLogger) prefixArgs(args []interface{}) string {
	return fmt.Sprint(append([]interface{}{l.Prefix + ": "}, args...)...)
}

func (l *PrefixedLogger) prefixMsg(message string) string {
	return fmt.Sprintf("%s: %s", l.Prefix, message)
}

func (l *PrefixedLogger) Info(args ...interface{}) {
	l.LoggerInstance.Info(l.prefixArgs(args))
}

func (l *PrefixedLogger) Infof(message string, args ...interface{}) {
	l.LoggerInstance.Info(l.prefixMsg(fmt.Sprintf(message, args...)))
}

func (l *PrefixedLogger) Debug(args ...interface{}) {
	l.LoggerInstance.Debug(l.prefixArgs(args))
}

func (l *PrefixedLogger) Debugf(message string, args ...interface{}) {
	l.LoggerInstance.Debug(l.prefixMsg(fmt.Sprintf(message, args...)))
}

func (l *PrefixedLogger) Error(args ...interface{}) {
	l.LoggerInstance.Error(l.prefixArgs(args))
}

func (l *PrefixedLogger) Errorf(message string, args ...interface{}) {
	l.LoggerInstance.Error(l.prefixMsg(fmt.Sprintf(message, args...)))
}

func (l *PrefixedLogger) Warn(args ...interface{}) {
	l.LoggerInstance.Warn(l.prefixArgs(args))
}

func (l *PrefixedLogger) Warnf(message string, args ...interface{}) {
	l.LoggerInstance.Warn(l.prefixMsg(fmt.Sprintf(message, args...)))
}

func (l *PrefixedLogger) Fatal(args ...interface{}) {
	l.LoggerInstance.Error(l.prefixArgs(args))
}

func (l *PrefixedLogger) Fatalf(message string, args ...interface{}) {
	l.LoggerInstance.Error(l.prefixMsg(fmt.Sprintf(message, args...)))
}

func (l *PrefixedLogger) Panic(args ...interface{}) {
	msg := l.prefixArgs(args)
	l.LoggerInstance.Error(msg)
	panic(msg)
}

func (l *PrefixedLogger) Panicf(message string, args ...interface{}) {
	msg := l.prefixMsg(fmt.Sprintf(message, args...))
	l.LoggerInstance.Error(msg)
	panic(msg)
}

// fields
func (l *PrefixedLogger) InfoWithFields(fields Fields, args ...interface{}) {
	l.LoggerInstance.With(toAttrs(fields)...).Info(l.prefixArgs(args))
}

func (l *PrefixedLogger) InfoWithFieldsf(fields Fields, message string, args ...interface{}) {
	l.LoggerInstance.With(toAttrs(fields)...).Info(l.prefixMsg(fmt.Sprintf(message, args...)))
}

func (l *PrefixedLogger) DebugWithFields(fields Fields, args ...interface{}) {
	l.LoggerInstance.With(toAttrs(fields)...).Debug(l.prefixArgs(args))
}

func (l *PrefixedLogger) DebugWithFieldsf(fields Fields, message string, args ...interface{}) {
	l.LoggerInstance.With(toAttrs(fields)...).Debug(l.prefixMsg(fmt.Sprintf(message, args...)))
}

func (l *PrefixedLogger) ErrorWithFields(fields Fields, args ...interface{}) {
	l.LoggerInstance.With(toAttrs(fields)...).Error(l.prefixArgs(args))
}

func (l *PrefixedLogger) ErrorWithFieldsf(fields Fields, message string, args ...interface{}) {
	l.LoggerInstance.With(toAttrs(fields)...).Error(l.prefixMsg(fmt.Sprintf(message, args...)))
}

func (l *PrefixedLogger) WarnWithFields(fields Fields, args ...interface{}) {
	l.LoggerInstance.With(toAttrs(fields)...).Warn(l.prefixArgs(args))
}

func (l *PrefixedLogger) WarnWithFieldsf(fields Fields, message string, args ...interface{}) {
	l.LoggerInstance.With(toAttrs(fields)...).Warn(l.prefixMsg(fmt.Sprintf(message, args...)))
}

func (l *PrefixedLogger) FatalWithFields(fields Fields, args ...interface{}) {
	l.LoggerInstance.With(toAttrs(fields)...).Error(l.prefixArgs(args))
}

func (l *PrefixedLogger) FatalWithFieldsf(fields Fields, message string, args ...interface{}) {
	l.LoggerInstance.With(toAttrs(fields)...).Error(l.prefixMsg(fmt.Sprintf(message, args...)))
}

func (l *PrefixedLogger) PanicWithFields(fields Fields, args ...interface{}) {
	msg := l.prefixArgs(args)
	l.LoggerInstance.With(toAttrs(fields)...).Error(msg)
	panic(msg)
}

func (l *PrefixedLogger) PanicWithFieldsf(fields Fields, message string, args ...interface{}) {
	msg := l.prefixMsg(fmt.Sprintf(message, args...))
	l.LoggerInstance.With(toAttrs(fields)...).Error(msg)
	panic(msg)
}

// error wrapping
func (l *PrefixedLogger) PrefixError(err error, msg string) error {
	return errors.Wrapf(err, "%s: %s", l.Prefix, msg)
}

func (l *PrefixedLogger) WrapError(err error) error {
	return errors.Wrapf(err, "%s: ", l.Prefix)
}
