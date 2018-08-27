package options

import (
	"github.com/spf13/pflag"

	"github.com/fengyfei/github-detector/cmd/github-detector/app/config"
	"github.com/fengyfei/github-detector/pkg/filetool"
)

// GitHubDetectorOptions is the main context object for the detector.
type GitHubDetectorOptions struct {
	ConfigurationPath string
	MirrorPath        string
	StorePath         string
}

// NewGitHubDetectorOptions creates a new nil GitHubDetectorOptions.
func NewGitHubDetectorOptions() *GitHubDetectorOptions {
	return &GitHubDetectorOptions{}
}

// AddFlags adds flags for a specific GitHubDetectorOptions to the specified FlagSet.
func (s *GitHubDetectorOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ConfigurationPath, "config", "config.json", "Specify configuration file path.")
	fs.StringVar(&s.MirrorPath, "mirror", "mirror.yaml", "Specify mirror repositories' file path.")
	fs.StringVar(&s.StorePath, "store", ".", "Specify the path to store files.")
}

// Config return a detector config objective.
func (s *GitHubDetectorOptions) Config() (*config.GitHubDetectorConfiguration, error) {
	cp, err := filetool.Abs(s.ConfigurationPath)
	if err != nil {
		return nil, err
	}

	mp, err := filetool.Abs(s.MirrorPath)
	if err != nil {
		return nil, err
	}

	sp, err := filetool.Abs(s.StorePath)
	if err != nil {
		return nil, err
	}

	return &config.GitHubDetectorConfiguration{
		ConfigurationPath: cp,
		MirrorPath:        mp,
		StorePath:         sp,
	}, nil
}
