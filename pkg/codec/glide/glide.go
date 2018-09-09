/*
 * Revision History:
 *     Initial: 2018/07/27        Li Zebang
 */

package glide

import (
	"time"

	"gopkg.in/yaml.v2"

	"github.com/TechCatsLab/github-detector/pkg/codec/conf"
	"github.com/TechCatsLab/github-detector/pkg/codec/lock"
)

// The Config struct from glide.
//
// https://raw.githubusercontent.com/Masterminds/glide/master/cfg/config.go
//
// Config is the top-level configuration object.
type Config struct {

	// Name is the name of the package or application.
	Name string `yaml:"package"`

	// Description is a short description for a package, application, or library.
	// This description is similar but different to a Go package description as
	// it is for marketing and presentation purposes rather than technical ones.
	Description string `json:"description,omitempty"`

	// Home is a url to a website for the package.
	Home string `yaml:"homepage,omitempty"`

	// License provides either a SPDX license or a path to a file containing
	// the license. For more information on SPDX see http://spdx.org/licenses/.
	// When more than one license an SPDX expression can be used.
	License string `yaml:"license,omitempty"`

	// Owners is an array of owners for a project. See the Owner type for
	// more detail. These can be one or more people, companies, or other
	// organizations.
	Owners Owners `yaml:"owners,omitempty"`

	// Ignore contains a list of packages to ignore fetching. This is useful
	// when walking the package tree (including packages of packages) to list
	// those to skip.
	Ignore []string `yaml:"ignore,omitempty"`

	// Exclude contains a list of directories in the local application to
	// exclude from scanning for dependencies.
	Exclude []string `yaml:"excludeDirs,omitempty"`

	// Imports contains a list of all non-development imports for a project. For
	// more detail on how these are captured see the Dependency type.
	Imports Dependencies `yaml:"import"`

	// DevImports contains the test or other development imports for a project.
	// See the Dependency type for more details on how this is recorded.
	DevImports Dependencies `yaml:"testImport,omitempty"`
}

// Owners is a list of owners for a project.
type Owners []*Owner

// Owner describes an owner of a package. This can be a person, company, or
// other organization. This is useful if someone needs to contact the
// owner of a package to address things like a security issue.
type Owner struct {

	// Name describes the name of an organization.
	Name string `yaml:"name,omitempty"`

	// Email is an email address to reach the owner at.
	Email string `yaml:"email,omitempty"`

	// Home is a url to a website for the owner.
	Home string `yaml:"homepage,omitempty"`
}

// Dependencies is a collection of Dependency
type Dependencies []*Dependency

// Dependency describes a package that the present package depends upon.
type Dependency struct {
	Name        string   `yaml:"package"`
	Reference   string   `yaml:"version,omitempty"`
	Pin         string   `yaml:"-"`
	Repository  string   `yaml:"repo,omitempty"`
	VcsType     string   `yaml:"vcs,omitempty"`
	Subpackages []string `yaml:"subpackages,omitempty"`
	Arch        []string `yaml:"arch,omitempty"`
	Os          []string `yaml:"os,omitempty"`
}

// The Lockfile struct from glide.
//
// https://raw.githubusercontent.com/Masterminds/glide/master/cfg/lock.go
//
// Lockfile represents a glide.lock file.
type Lockfile struct {
	Hash       string    `yaml:"hash"`
	Updated    time.Time `yaml:"updated"`
	Imports    Locks     `yaml:"imports"`
	DevImports Locks     `yaml:"testImports"`
}

// Locks is a slice of locked dependencies.
type Locks []*Lock

// Lock represents an individual locked dependency.
type Lock struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Repository  string   `yaml:"repo,omitempty"`
	VcsType     string   `yaml:"vcs,omitempty"`
	Subpackages []string `yaml:"subpackages,omitempty"`
	Arch        []string `yaml:"arch,omitempty"`
	Os          []string `yaml:"os,omitempty"`
}

type glideParser struct{}

// GlidePraser -
var GlidePraser glideParser

// ParseConfig -
func (gp *glideParser) ParseConfig(in []byte) (*Config, error) {
	var config Config
	return &config, yaml.Unmarshal(in, &config)
}

// ParseLockfile -
func (gp *glideParser) ParseLockfile(in []byte) (*Lockfile, error) {
	var lockfile Lockfile
	return &lockfile, yaml.Unmarshal(in, &lockfile)
}

// Deps -
func (c *Config) Deps() []string {
	deps := make([]string, len(c.Imports))
	for index := range c.Imports {
		deps[index] = c.Imports[index].Name
	}
	return deps
}

// Repos -
func (l *Lockfile) Repos() []string {
	repos := make([]string, len(l.Imports))
	for index := range l.Imports {
		repos[index] = l.Imports[index].Name
	}
	return repos
}

// ParseConfFile -
func (gp *glideParser) ParseConfFile(in []byte) (conf.File, error) {
	return gp.ParseConfig(in)
}

// ParseLockFile -
func (gp *glideParser) ParseLockFile(in []byte) (lock.File, error) {
	return gp.ParseLockfile(in)
}
