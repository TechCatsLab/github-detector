/*
 * Revision History:
 *     Initial: 2018/08/02        Li Zebang
 */

package mirror

import (
	store "github.com/fengyfei/github-detector/pkg/store/file/yaml"
)

// Mirror -
type Mirror struct {
	Package string `yaml:"package"`
	Repo    string `yaml:"repo"`
	VCS     string `yaml:"vcs"`
}

// LoadMirrors -
func LoadMirrors(path string) (*store.File, error) {
	return store.Open(path)
}

// ExportMirrors -
func ExportMirrors(f *store.File) error {
	return f.SaveAndClose()
}
