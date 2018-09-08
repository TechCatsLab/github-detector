/*
 * Revision History:
 *     Initial: 2018/08/02        Li Zebang
 */

package mirror

import (
	"github.com/TechCatsLab/github-detector/pkg/sync"
)

// Mirror -
type Mirror struct {
	Package string `json:"package"`
	Repo    string `json:"repo"`
	VCS     string `json:"vcs"`
}

// LoadMirrors -
func LoadMirrors(path string) (*sync.Map, error) {
	return sync.Open(path)
}

// ExportMirrors -
func ExportMirrors(m *sync.Map, path string) error {
	return m.Save(path)
}
