/*
 * Revision History:
 *     Initial: 2018/08/02        Li Zebang
 */

package app

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/scheduler"

	"github.com/fengyfei/github-detector/cmd/github-detector/app/config"
	"github.com/fengyfei/github-detector/cmd/github-detector/app/mirror"
	"github.com/fengyfei/github-detector/cmd/github-detector/app/options"
	"github.com/fengyfei/github-detector/pkg/filetool"
	pool "github.com/fengyfei/github-detector/pkg/github"
	store "github.com/fengyfei/github-detector/pkg/store/file/yaml"
	"github.com/fengyfei/github-detector/pkg/version"
	"github.com/fengyfei/github-detector/pkg/version/verflag"
)

// NewDetector creates a *cobra.Command object with default parameters.
func NewDetector() *cobra.Command {
	s := options.NewGitHubDetectorOptions()

	cmd := &cobra.Command{
		Use: "github-detector",
		Long: `The GitHub detector is a git repositories' detector tool.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			verflag.GetVersion()

			c, err := s.Config()
			if err != nil {
				logrus.Fatalf("%v", err)
			}

			if err := Run(c); err != nil {
				logrus.Fatalf("%v", err)
			}
		},
	}
	s.AddFlags(cmd.Flags())

	return cmd
}

// Run runs the detector.
func Run(c *config.GitHubDetectorConfiguration) error {
	logrus.Infof("Version: %+v", version.Get())

	var err error
	c.Configuration, err = config.LoadConfiguration(c.ConfigurationPath)
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	c.Mirrors, err = mirror.LoadMirrors(c.MirrorPath)
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	c.SPool = scheduler.New(len(c.Configuration.Tokens), len(c.Configuration.Tokens))
	c.GPool = pool.NewPool(c.Configuration.Tokens...)
	err = os.MkdirAll(c.StorePath+"/repos/", 0755)
	if err != nil {
		logrus.Fatalf("%v", err)
	}
	err = os.MkdirAll(c.StorePath+"/cache/", 0755)
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	return run(c)
}

func run(c *config.GitHubDetectorConfiguration) error {
	stc := NewSearchTaskContext(&SearchTaskInfo{
		Dir:      c.StorePath,
		Language: c.Configuration.Language,
		Pushed:   c.Configuration.Pushed,
		Min:      c.Configuration.Min,
		Max:      c.Configuration.Max,
		GPool:    c.GPool,
	})
	srt := NewRetryTask(stc, scheduler.TaskFunc(SearchTaskFunc))
	c.SPool.Schedule(NewRetryTaskContext(&RetryTaskInfo{
		Times: 3,
		SPool: c.SPool,
	}), srt)
	c.SPool.Wait()
	logrus.Infof("SEARCH TASK FINISHED")

	type r struct {
		FullName string `json:"full_name"`
	}
	rs, err := func() ([]r, error) {
		rs := make([]r, 0)
		f, err := filetool.Open(c.StorePath+"/repos.json", filetool.RDONLY, 0644)
		if err != nil {
			return nil, err
		}
		return rs, filetool.NewDecoder(f).Decode(&rs)
	}()
	if err != nil {
		logrus.Errorf("Opening repos.json failed, %v", err)
		return err
	}

	infof, err := store.Open(c.StorePath + "/info.yaml")
	if err != nil {
		logrus.Errorf("Opening info.yaml failed, %v", err)
		return err
	}
	defer infof.SaveAndClose()

	for index := range rs {
		i := strings.Index(rs[index].FullName, "/")
		if i == -1 || i+2 > len(rs[index].FullName) {
			logrus.Errorf("invalid full name: %s", rs[index])
			continue
		}
		ltc := NewListTaskContext(&ListTaskInfo{
			Dir:   c.StorePath + "/cache/",
			URL:   "github.com/" + rs[index].FullName,
			Owner: rs[index].FullName[:i],
			Repo:  rs[index].FullName[i+1:],
			Info:  infof,
			GPool: c.GPool,
		})
		lrt := NewRetryTask(ltc, scheduler.TaskFunc(ListTaskFunc))
		c.SPool.Schedule(NewRetryTaskContext(&RetryTaskInfo{
			Times: 3,
			SPool: c.SPool,
		}), lrt)
	}
	c.SPool.Wait()
	logrus.Info("LIST TASK FINISHED")

	itc := NewIndexTaskContext(&IndexTaskInfo{
		CacheDir: c.StorePath + "/cache/",
		ReposDir: c.StorePath + "/repos/",
		Info:     infof,
	})
	irt := NewRetryTask(itc, scheduler.TaskFunc(IndexTaskFunc))
	c.SPool.Schedule(NewRetryTaskContext(&RetryTaskInfo{
		Times: 3,
		SPool: c.SPool,
	}), irt)
	c.SPool.Wait()
	logrus.Info("INDEX TASK FINISHED")

	return nil
}
