/*
 * Revision History:
 *     Initial: 2018/07/27        Li Zebang
 */

package godep

import (
	"encoding/json"

	"github.com/fengyfei/github-detector/pkg/codec/conf"
	"github.com/fengyfei/github-detector/pkg/codec/lock"
)

// The Godeps struct from godep.
//
// https://raw.githubusercontent.com/tools/godep/master/godepfile.go
//
// Godeps describes what a package needs to be rebuilt reproducibly.
// It's the same information stored in file Godeps.
type Godeps struct {
	ImportPath   string       `json:"ImportPath"`
	GoVersion    string       `json:"GoVersion"`
	GodepVersion string       `json:"GodepVersion"`
	Packages     []string     `json:"Packages,omitempty"` // Arguments to save, if any.
	Deps         []Dependency `json:"Deps"`
}

// The Dependency struct from godep.
//
// https://raw.githubusercontent.com/tools/godep/master/dep.go
//
// A Dependency is a specific revision of a package.
type Dependency struct {
	ImportPath string `json:"ImportPath"`
	Comment    string `json:"Comment,omitempty"` // Description of commit, if present.
	Rev        string `json:"Rev"`               // VCS-specific commit ID.
}

type godepParser struct{}

// GodepParser -
var GodepParser godepParser

// ParseGodeps -
func (gp *godepParser) ParseGodeps(in []byte) (*Godeps, error) {
	var godeps Godeps
	return &godeps, json.Unmarshal(in, &godeps)
}

// Repos -
func (g *Godeps) Repos() []string {
	repos := make([]string, len(g.Deps))
	for index := range g.Deps {
		repos[index] = g.Deps[index].ImportPath
	}
	return repos
}

// ParseConfFile -
func (gp *godepParser) ParseConfFile(in []byte) (conf.File, error) {
	return nil, nil
}

// ParseLockFile -
func (gp *godepParser) ParseLockFile(in []byte) (lock.File, error) {
	return gp.ParseGodeps(in)
}
