/*
 * Revision History:
 *     Initial: 2018/08/01        Li Zebang
 */

package filetool

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Option -
type Option int

const (
	// IGNORE -
	IGNORE Option = iota
	// RDONLY -
	RDONLY
	// WRONLY -
	WRONLY
	// RDWR -
	RDWR
	// APPEND -
	APPEND
	// TRUNC -
	TRUNC
)

// Open -
func Open(name string, opt Option, perm os.FileMode) (*os.File, error) {
	if perm == 0 {
		perm = 0644
	}
	switch opt {
	case RDONLY:
		return os.OpenFile(name, os.O_CREATE|os.O_RDONLY, perm)
	case WRONLY:
		return os.OpenFile(name, os.O_CREATE|os.O_WRONLY, perm)
	case RDWR:
		return os.OpenFile(name, os.O_CREATE|os.O_RDWR, perm)
	case APPEND:
		return os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_APPEND, perm)
	case TRUNC:
		return os.OpenFile(name, os.O_CREATE|os.O_RDWR|os.O_TRUNC, perm)
	default:
		return nil, errors.New("unsupported option")
	}
}

// Abs -
func Abs(name string) (string, error) {
	return filepath.Abs(name)
}

// Ext -
func Ext(name string) string {
	return filepath.Ext(name)
}

// NewEncoder -
func NewEncoder(f *os.File) *json.Encoder {
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	return encoder
}

// NewDecoder -
func NewDecoder(f *os.File) *json.Decoder {
	return json.NewDecoder(f)
}
