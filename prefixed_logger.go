package log

import (
	"fmt"
	"io"

	"github.com/neighborly/go-errors"
	"github.com/sirupsen/logrus"
)

type PrefixedLogger struct {
	Prefix string
}

func NewPrefixedLogger(prefix string) PrefixedLogger {
	return PrefixedLogger{
		Prefix: prefix,
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
	logger.Info(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Infof(message string, args ...interface{}) {
	logger.Infof(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Debug(args ...interface{}) {
	logger.Debug(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Debugf(message string, args ...interface{}) {
	logger.Debugf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Error(args ...interface{}) {
	logger.Error(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Errorf(message string, args ...interface{}) {
	logger.Errorf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Warn(args ...interface{}) {
	logger.Warn(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Warnf(message string, args ...interface{}) {
	logger.Warnf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Fatal(args ...interface{}) {
	logger.Fatal(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Fatalf(message string, args ...interface{}) {
	logger.Fatalf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Panic(args ...interface{}) {
	logger.Panic(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) Panicf(message string, args ...interface{}) {
	logger.Panicf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) Writer() *io.PipeWriter {
	return logger.Writer()
}

// fields

func (l *PrefixedLogger) InfoWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Info(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) InfoWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Infof(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) DebugWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Debug(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) DebugWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Debugf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) ErrorWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Error(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) ErrorWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Errorf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) WarnWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Warn(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) WarnWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Warnf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) FatalWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Fatal(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) FatalWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Fatalf(l.prefixMsg(message), args...)
}

func (l *PrefixedLogger) PanicWithFields(fields Fields, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Panic(l.prefixArgs(args)...)
}

func (l *PrefixedLogger) PanicWithFieldsf(fields Fields, message string, args ...interface{}) {
	logger.WithFields(logrus.Fields(fields)).Panicf(l.prefixMsg(message), args...)
}

// error wrapping

func (l *PrefixedLogger) PrefixError(err error, msg string) error {
	return errors.Wrapf(err, "%s: %s", l.Prefix, msg)
}

func (l *PrefixedLogger) WrapError(err error) error {
	return errors.Wrapf(err, "%s: ", l.Prefix)
}
