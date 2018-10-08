/*
 * Revision History:
 *     Initial: 2018/08/12        Feng Yifei
 */

package logging

// Logger represents a general purpose logging functions.
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Option is used for setting up a Logger.
type Option func(Logger) error
