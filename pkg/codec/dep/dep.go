/*
 * Revision History:
 *     Initial: 2018/07/27        Li Zebang
 */

package dep

import (
	"github.com/pelletier/go-toml"

	"github.com/fengyfei/github-detector/pkg/codec/conf"
	"github.com/fengyfei/github-detector/pkg/codec/lock"
)

// ManifestName is the manifest file name used by dep.
const ManifestName = "Gopkg.toml"

// The Manifest struct from dep.
//
// https://raw.githubusercontent.com/golang/dep/master/manifest.go
//
// Manifest represents a Gopkg.toml file.
type Manifest struct {
	Constraints  []Project    `toml:"constraint,omitempty" json:"constraint,omitempty"`
	Overrides    []Project    `toml:"override,omitempty" json:"override,omitempty"`
	Ignored      []string     `toml:"ignored,omitempty" json:"ignored,omitempty"`
	Required     []string     `toml:"required,omitempty" json:"required,omitempty"`
	NoVerify     []string     `toml:"noverify,omitempty" json:"noverify,omitempty"`
	PruneOptions PruneOptions `toml:"prune,omitempty" json:"prune,omitempty"`
}

// Project -
type Project struct {
	Name     string `toml:"name" json:"name"`
	Branch   string `toml:"branch,omitempty" json:"branch,omitempty"`
	Revision string `toml:"revision,omitempty" json:"revision,omitempty"`
	Version  string `toml:"version,omitempty" json:"version,omitempty"`
	Source   string `toml:"source,omitempty" json:"source,omitempty"`
}

// PruneOptions -
type PruneOptions struct {
	UnusedPackages bool `toml:"unused-packages,omitempty" json:"unused-packages,omitempty"`
	NonGoFiles     bool `toml:"non-go,omitempty" json:"non-go,omitempty"`
	GoTests        bool `toml:"go-tests,omitempty" json:"go-tests,omitempty"`

	//Projects []map[string]interface{} `toml:"project,omitempty"`
	Projects []map[string]interface{}
}

// LockName is the lock file name used by dep.
const LockName = "Gopkg.lock"

// The Lock struct from dep.
//
// https://raw.githubusercontent.com/golang/dep/master/lock.go
//
// Lock represents a Gopkg.lock file.
type Lock struct {
	SolveMeta Meta            `toml:"solve-meta" json:"solve-meta"`
	Projects  []LockedProject `toml:"projects" json:"projects"`
}

// Meta -
type Meta struct {
	AnalyzerName    string   `toml:"analyzer-name" json:"analyzer-name"`
	AnalyzerVersion int      `toml:"analyzer-version" json:"analyzer-version"`
	SolverName      string   `toml:"solver-name" json:"solver-name"`
	SolverVersion   int      `toml:"solver-version" json:"solver-version"`
	InputImports    []string `toml:"input-imports" json:"input-imports"`
}

// LockedProject -
type LockedProject struct {
	Name      string   `toml:"name" json:"name"`
	Branch    string   `toml:"branch,omitempty" json:"branch,omitempty"`
	Revision  string   `toml:"revision" json:"revision"`
	Version   string   `toml:"version,omitempty" json:"version,omitempty"`
	Source    string   `toml:"source,omitempty" json:"source,omitempty"`
	Packages  []string `toml:"packages" json:"packages"`
	PruneOpts string   `toml:"pruneopts" json:"pruneopts"`
	Digest    string   `toml:"digest" json:"digest"`
}

type depParser struct{}

// DepParser -
var DepParser depParser

// ParseToml -
func (dp *depParser) ParseToml(in []byte) (*Manifest, error) {
	var manifest Manifest
	return &manifest, toml.Unmarshal(in, &manifest)
}

// ParseLock -
func (dp *depParser) ParseLock(in []byte) (*Lock, error) {
	var lock Lock
	return &lock, toml.Unmarshal(in, &lock)
}

// Deps -
func (m *Manifest) Deps() []string {
	deps := make([]string, len(m.Constraints))
	for index := range m.Constraints {
		deps[index] = m.Constraints[index].Name
	}
	return deps
}

// Repos -
func (l *Lock) Repos() []string {
	repos := make([]string, len(l.Projects))
	for index := range l.Projects {
		repos[index] = l.Projects[index].Name
	}
	return repos
}

// ParseConfFile -
func (dp *depParser) ParseConfFile(in []byte) (conf.File, error) {
	return dp.ParseToml(in)
}

// ParseLockFile -
func (dp *depParser) ParseLockFile(in []byte) (lock.File, error) {
	return dp.ParseLock(in)
}
