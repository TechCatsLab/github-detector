package version

import (
	"fmt"
	"runtime"
)

const (
	version = "0.0.1"
)

type Info struct {
	Version   string `json:"version"`
	GoVersion string `json:"goVersion"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
}

// Print prints version information.
func (i *Info) Print() {
	fmt.Printf(`Version Information:
%-9s %s
%-9s %s
%-9s %s
%-9s %s
`,
		"Version", i.Version,
		"GoVersion", i.GoVersion,
		"Compiler", i.Compiler,
		"Platform", i.Platform,
	)
}

// Get returns the version information.
func Get() *Info {
	return &Info{
		Version:   version,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
