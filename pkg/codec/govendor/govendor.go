/*
 * Revision History:
 *     Initial: 2018/07/27        Li Zebang
 */

package govendor

import (
	"encoding/json"

	"github.com/TechCatsLab/github-detector/pkg/codec/conf"
	"github.com/TechCatsLab/github-detector/pkg/codec/lock"
)

// Name of the vendor file.
const Name = "vendor.json"

// The File struct from govendor.
//
// https://raw.githubusercontent.com/kardianos/govendor/master/vendorfile/file.go
//
// File is the structure of the vendor file.
type File struct {
	RootPath string `json:"rootPath"` // Import path of vendor folder

	Comment string `json:"comment"`

	Ignore string `json:"ignore"`

	Package []*Package `json:"package"`
}

// Package represents each package.
type Package struct {
	// See the vendor spec for definitions.
	Origin       string `json:"origin,omitempty"`
	Path         string `json:"path"`
	Tree         bool   `json:"tree,omitempty"`
	Revision     string `json:"revision"`
	RevisionTime string `json:"revisionTime"`
	Version      string `json:"version,omitempty"`
	VersionExact string `json:"versionExact,omitempty"`
	ChecksumSHA1 string `json:"checksumSHA1"`
	Comment      string `json:"comment,omitempty"`

	ImportPath string `json:"importpath,omitempty"`
}

type govendorParser struct{}

// GovendorParser -
var GovendorParser govendorParser

// ParseFile -
func (gp *govendorParser) ParseFile(in []byte) (*File, error) {
	var file File
	return &file, json.Unmarshal(in, &file)
}

// Repos -
func (f *File) Repos() []string {
	repos := make([]string, len(f.Package))
	for index := range f.Package {
		if f.Package[index].ImportPath != "" {
			repos[index] = f.Package[index].ImportPath
		} else {
			repos[index] = f.Package[index].Path
		}
	}
	return repos
}

// ParseConfFile -
func (gp *govendorParser) ParseConfFile(in []byte) (conf.File, error) {
	return nil, nil
}

// ParseLockFile -
func (gp *govendorParser) ParseLockFile(in []byte) (lock.File, error) {
	return gp.ParseFile(in)
}
