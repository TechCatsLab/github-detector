/*
 * Revision History:
 *     Initial: 2018/08/02        Li Zebang
 */

package config

import (
	"errors"
	"time"

	"github.com/TechCatsLab/scheduler"

	"github.com/TechCatsLab/github-detector/pkg/filetool"
	"github.com/TechCatsLab/github-detector/pkg/github"
	"github.com/TechCatsLab/github-detector/pkg/sync"
)

// GitHubDetectorConfiguration is the main context object for the detector.
type GitHubDetectorConfiguration struct {
	ConfigurationPath string
	MirrorPath        string
	StorePath         StorePath

	Configuration *Configuration
	Mirrors       *sync.Map

	GPool github.Pool
	SPool *scheduler.Pool
}

// StorePath -
type StorePath struct {
	Root  string
	Repo  string
	Cache string
	Repos string
}

// Configuration -
type Configuration struct {
	Interval time.Duration
	Language string
	Pushed   time.Duration
	Min      int
	Max      int
	Tokens   []*github.Token
}

type configuration struct {
	Interval string `json:"interval"`
	Language string `json:"language"`
	Pushed   string `json:"pushed"`
	Stars    struct {
		Min int `json:"min,omitempty"`
		Max int `json:"max"`
	} `json:"stars"`
	Tokens []*github.Token `json:"tokens"`
}

// LoadConfiguration -
func LoadConfiguration(path string) (*Configuration, error) {
	file, err := filetool.Open(path, filetool.RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var c configuration
	err = filetool.NewDecoder(file).Decode(&c)
	if err != nil {
		return nil, err
	}

	if c.Stars.Min < 0 || c.Stars.Max < 0 {
		return nil, errors.New("stars number in the configuration file is invalid")
	}

	if c.Stars.Min > c.Stars.Max {
		c.Stars.Min = c.Stars.Max
	}

	interva, err := StringToTime(c.Interval, " ")
	if err != nil {
		return nil, err
	}

	pushed, err := StringToTime(c.Pushed, " ")
	if err != nil {
		return nil, err
	}

	return &Configuration{
		Interval: interva,
		Language: c.Language,
		Pushed:   pushed,
		Min:      c.Stars.Min,
		Max:      c.Stars.Max,
		Tokens:   c.Tokens,
	}, nil
}

// ExportConfiguration -
func ExportConfiguration(c *Configuration, path string) error {
	var config = &configuration{
		Language: c.Language,
		Tokens:   c.Tokens,
	}
	config.Stars.Min = c.Min
	config.Stars.Max = c.Max

	interval, err := TimeToString(c.Interval, " ")
	if err != nil {
		return err
	}
	config.Interval = interval

	pushed, err := TimeToString(c.Pushed, " ")
	if err != nil {
		return err
	}
	config.Pushed = pushed

	file, err := filetool.Open(path, filetool.TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return filetool.NewEncoder(file).Encode(config)
}
