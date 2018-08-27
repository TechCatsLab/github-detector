/*
 * Revision History:
 *     Initial: 2018/08/02        Li Zebang
 */

package config

import (
	"errors"
	"time"

	"github.com/TechCatsLab/scheduler"

	"github.com/fengyfei/github-detector/pkg/filetool"
	"github.com/fengyfei/github-detector/pkg/github"
	store "github.com/fengyfei/github-detector/pkg/store/file/yaml"
)

// GitHubDetectorConfiguration is the main context object for the detector.
type GitHubDetectorConfiguration struct {
	ConfigurationPath string
	MirrorPath        string
	StorePath         string

	Configuration *Configuration
	Mirrors       *store.File

	GPool github.Pool
	SPool *scheduler.Pool
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
	f, err := filetool.Open(path, filetool.RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c configuration
	err = filetool.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, err
	}

	if c.Stars.Min < 0 || c.Stars.Max < 0 {
		return nil, errors.New("stars number in the configuration file is invalid")
	}

	if c.Stars.Min > c.Stars.Max {
		c.Stars.Min = c.Stars.Max
	}

	i, err := StringToTime(c.Interval, " ")
	if err != nil {
		return nil, err
	}

	p, err := StringToTime(c.Pushed, " ")
	if err != nil {
		return nil, err
	}

	return &Configuration{
		Interval: i,
		Language: c.Language,
		Pushed:   p,
		Min:      c.Stars.Min,
		Max:      c.Stars.Max,
		Tokens:   c.Tokens,
	}, nil
}

// ExportLoadConfiguration -
func ExportLoadConfiguration(c *Configuration, path string) error {
	var config = &configuration{
		Language: c.Language,
		Tokens:   c.Tokens,
	}
	config.Stars.Min = c.Min
	config.Stars.Max = c.Max

	i, err := TimeToString(c.Interval, " ")
	if err != nil {
		return err
	}
	config.Interval = i

	p, err := TimeToString(c.Pushed, " ")
	if err != nil {
		return err
	}
	config.Pushed = p

	f, err := filetool.Open(path, filetool.TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return filetool.NewEncoder(f).Encode(config)
}
