/*
 * Revision History:
 *     Initial: 2018/07/27        Li Zebang
 */

package codec

import (
	"github.com/TechCatsLab/github-detector/pkg/codec/conf"
	"github.com/TechCatsLab/github-detector/pkg/codec/dep"
	"github.com/TechCatsLab/github-detector/pkg/codec/glide"
	"github.com/TechCatsLab/github-detector/pkg/codec/godep"
	"github.com/TechCatsLab/github-detector/pkg/codec/govendor"
	"github.com/TechCatsLab/github-detector/pkg/codec/lock"
)

// Codec -
type Codec interface {
	ParseConfFile([]byte) (conf.File, error)
	ParseLockFile([]byte) (lock.File, error)
}

// Dep -
func Dep() Codec {
	return &dep.DepParser
}

// Gilde -
func Gilde() Codec {
	return &glide.GlidePraser
}

// Godep -
func Godep() Codec {
	return &godep.GodepParser
}

// Govendor -
func Govendor() Codec {
	return &govendor.GovendorParser
}
