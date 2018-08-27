/*
 * Revision History:
 *     Initial: 2018/08/12        Feng Yifei
 */

package logrus

import (
	"github.com/TechCatsLab/logging"
	log "github.com/sirupsen/logrus"
)

var (
	_ logging.Logger = &logger{}
)

type logger struct {
	lg *log.Logger
}

// New creates a logger instance.
func New(opts ...logging.Option) logging.Logger {
	l := &logger{
		lg: log.New(),
	}

	for _, fn := range opts {
		fn(l)
	}

	return l
}

// Debug - Logger implementation.
func (l *logger) Debug(args ...interface{}) {
	l.lg.Debug(args...)
}

// Info - Logger implementation.
func (l *logger) Info(args ...interface{}) {
	l.lg.Info(args...)
}

// Warn - Logger implementation.
func (l *logger) Warn(args ...interface{}) {
	l.lg.Warn(args...)
}

// Error - Logger implementation.
func (l *logger) Error(args ...interface{}) {
	l.lg.Error(args...)
}

// Fatal - Logger implementation.
func (l *logger) Fatal(args ...interface{}) {
	l.lg.Fatal(args...)
}

// Debugf - Logger implementation.
func (l *logger) Debugf(format string, args ...interface{}) {
	l.lg.Debugf(format, args...)
}

// Infof - Logger implementation.
func (l *logger) Infof(format string, args ...interface{}) {
	l.lg.Infof(format, args...)
}

// Warnf - Logger implementation.
func (l *logger) Warnf(format string, args ...interface{}) {
	l.lg.Warnf(format, args...)
}

// Errorf - Logger implementation.
func (l *logger) Errorf(format string, args ...interface{}) {
	l.lg.Errorf(format, args...)
}

// Fatalf - Logger implementation.
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.lg.Fatalf(format, args...)
}
