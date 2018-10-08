/*
 * Revision History:
 *     Initial: 2018/08/12        Feng Yifei
 */

package logrus

var (
	defaultLogger *logger

	// Debug is a global method for showing debuging messages.
	Debug func(...interface{})

	// Info is a global method for showing debuging messages.
	Info func(...interface{})

	// Warn is a global method for showing debuging messages.
	Warn func(...interface{})

	// Error is a global method for showing debuging messages.
	Error func(...interface{})

	// Fatal is a global method for showing debuging messages.
	Fatal func(...interface{})

	// Debugf is a global method for showing debuging messages.
	Debugf func(string, ...interface{})

	// Infof is a global method for showing debuging messages.
	Infof func(string, ...interface{})

	// Warnf is a global method for showing debuging messages.
	Warnf func(string, ...interface{})

	// Errorf is a global method for showing debuging messages.
	Errorf func(string, ...interface{})

	// Fatalf is a global method for showing debuging messages.
	Fatalf func(string, ...interface{})
)

func init() {
	defaultLogger := New(OptSetInfoLevel, OptShowFileLine)

	Debug = defaultLogger.Debug
	Info = defaultLogger.Info
	Warn = defaultLogger.Warn
	Error = defaultLogger.Error
	Fatal = defaultLogger.Fatal

	Debugf = defaultLogger.Debugf
	Infof = defaultLogger.Infof
	Warnf = defaultLogger.Warnf
	Errorf = defaultLogger.Errorf
	Fatalf = defaultLogger.Fatalf
}
