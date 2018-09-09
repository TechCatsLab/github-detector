package verflag

import (
	"os"

	"github.com/spf13/pflag"

	"github.com/TechCatsLab/github-detector/pkg/version"
)

var (
	versionFlag = pflag.BoolP("version", "v", false, "Get version information and quit.")
)

// GetVersion prints the version information and exit if the -version flag was passed.
func GetVersion() {
	if *versionFlag {
		version.Get().Print()
		os.Exit(0)
	}
}
