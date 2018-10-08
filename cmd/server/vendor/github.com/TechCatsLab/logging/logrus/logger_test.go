/*
 * Revision History:
 *     Initial: 2018/08/12        Feng Yifei
 */

package logrus

import (
	"testing"
)

func Test_Logger(t *testing.T) {
	log := New(OptSetDebugLevel, OptShowFileLine)

	log.Debug("Debug")
	log.Info("Info")
	log.Warn("Warn")
	log.Error("Error")
}

func Test_Default(t *testing.T) {
	Debug("Debug")
	Info("Info")
	Warn("Warn")
	Error("Error")
}
