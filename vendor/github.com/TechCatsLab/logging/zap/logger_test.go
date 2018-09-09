/*
 * Revision History:
 *     Initial: 2018/08/19        Feng Yifei
 */
package zap

import (
	"testing"
)

func Test_Default(t *testing.T) {
	Debug("Debug")
	Info("Info")
	Warn("Warn")
	Error("Error")
}
