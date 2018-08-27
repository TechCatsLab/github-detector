/*
 * Revision History:
 *     Initial: 2018/08/12        Feng Yifei
 */

package hooks

import (
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// FileLineHook is used for show file & lines.
type FileLineHook struct {
}

// Levels is required by logrus.
func (hook FileLineHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire is required by logrus.
func (hook FileLineHook) Fire(entry *log.Entry) error {
	if pc, file, line, ok := runtime.Caller(8); ok {
		funcName := runtime.FuncForPC(pc).Name()
		entry.Data["file"] = path.Base(file)
		entry.Data["func"] = path.Base(funcName)
		entry.Data["line"] = line
	}
	return nil
}
