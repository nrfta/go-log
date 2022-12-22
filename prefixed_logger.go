package log

import (
	"fmt"
	"io"

	"github.com/neighborly/go-errors"
	"github.com/sirupsen/logrus"
)

type PrefixedLogger struct {
	Prefix         string
	LoggerInstance *logrus.Logger
}

func NewPrefixedLogger(prefix string, instance *logrus.Logger) PrefixedLogger {
	if instance == nil {
		instance = New(true, "info")
	}

	return PrefixedLogger{
		Prefix:         prefix,
		LoggerInstance: instance,
	}
}

func (l *PrefixedLogger) prefixArgs(args []interface{}) []interface{} {
	return append([]interface{}{l.Prefix + ": "}, args...)
}

func (l *PrefixedLogger) prefixMsg(message string) string {
	return fmt.Sprintf("%s: %s", l.Prefix, message)
}

// implement logger interface
func (l *PrefixedLogger) Info(args ...interface{}) {
	l.LoggerInstance.Info(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Infof(message string, args ...interface{}) {
	l.LoggerInstance.Infof(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Debug(args ...interface{}) {
	l.LoggerInstance.Debug(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Debugf(message string, args ...interface{}) {
	l.LoggerInstance.Debugf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Error(args ...interface{}) {
	l.LoggerInstance.Error(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Errorf(message string, args ...interface{}) {
	l.LoggerInstance.Errorf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Warn(args ...interface{}) {
	l.LoggerInstance.Warn(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Warnf(message string, args ...interface{}) {
	l.LoggerInstance.Warnf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Fatal(args ...interface{}) {
	l.LoggerInstance.Fatal(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Fatalf(message string, args ...interface{}) {
	l.LoggerInstance.Fatalf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Panic(args ...interface{}) {
	l.LoggerInstance.Panic(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Panicf(message string, args ...interface{}) {
	l.LoggerInstance.Panicf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Writer() *io.PipeWriter {
	return l.LoggerInstance.Writer()
}

// fields

func (l *PrefixedLogger) InfoWithFields(fields Fields, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Info(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) InfoWithFieldsf(fields Fields, message string, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Infof(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) DebugWithFields(fields Fields, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Debug(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) DebugWithFieldsf(fields Fields, message string, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Debugf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) ErrorWithFields(fields Fields, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Error(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) ErrorWithFieldsf(fields Fields, message string, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Errorf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) WarnWithFields(fields Fields, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Warn(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) WarnWithFieldsf(fields Fields, message string, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Warnf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) FatalWithFields(fields Fields, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Fatal(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) FatalWithFieldsf(fields Fields, message string, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Fatalf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) PanicWithFields(fields Fields, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Panic(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) PanicWithFieldsf(fields Fields, message string, args ...interface{}) {
	l.LoggerInstance.WithFields(logrus.Fields(fields)).Panicf(l.prefixMsg(message), args...)
}

// error wrapping

func (l *PrefixedLogger) PrefixError(err error, msg string) error {
	return errors.Wrapf(err, "%s: %s", l.Prefix, msg)
}

func (l *PrefixedLogger) WrapError(err error) error {
	return errors.Wrapf(err, "%s: ", l.Prefix)
}
