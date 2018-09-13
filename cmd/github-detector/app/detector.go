/*
 * Revision History:
 *     Initial: 2018/08/02        Li Zebang
 */

package app

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/scheduler"

	"github.com/TechCatsLab/github-detector/cmd/github-detector/app/config"
	"github.com/TechCatsLab/github-detector/cmd/github-detector/app/mirror"
	"github.com/TechCatsLab/github-detector/cmd/github-detector/app/options"
	"github.com/TechCatsLab/github-detector/pkg/filetool"
	pool "github.com/TechCatsLab/github-detector/pkg/github"
	"github.com/TechCatsLab/github-detector/pkg/sync"
	"github.com/TechCatsLab/github-detector/pkg/version"
	"github.com/TechCatsLab/github-detector/pkg/version/verflag"
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
	defer func() {
		if err != nil {
			logrus.Fatalf("%v", err)
		}
	}()

	c.Configuration, err = config.LoadConfiguration(c.ConfigurationPath)
	if err != nil {
		return err
	}
	c.Mirrors, err = mirror.LoadMirrors(c.MirrorPath)
	if err != nil {
		return err
	}
	c.SPool = scheduler.New(len(c.Configuration.Tokens), len(c.Configuration.Tokens))
	c.GPool = pool.NewPool(c.Configuration.Tokens...)

	var now time.Time
	for {
		now = time.Now()
		c.StorePath.Repo = c.StorePath.Root + RepoDir + now.Format("2006-01-02") + "/"
		c.StorePath.Repos = c.StorePath.Root + RepoDir + now.Format("2006-01-02") + "/" + ReposDir
		c.StorePath.Cache = c.StorePath.Root + RepoDir + now.Format("2006-01-02") + "/" + CacheDir
		err = os.MkdirAll(c.StorePath.Repo, 0755)
		if err != nil {
			return err
		}
		err = os.MkdirAll(c.StorePath.Repos, 0755)
		if err != nil {
			return err
		}
		err = os.MkdirAll(c.StorePath.Cache, 0755)
		if err != nil {
			return err
		}

		if err = run(c); err != nil {
			logrus.Errorf("Running failed, %v", err)
			os.RemoveAll(c.StorePath.Repo)
			<-time.After(time.Hour)
			continue
		}

		<-time.After(time.Until(now.Add(c.Configuration.Interval)))
	}
}

func run(c *config.GitHubDetectorConfiguration) error {
	stc := NewSearchTaskContext(&SearchTaskInfo{
		RepoDir:  c.StorePath.Repo,
		Language: c.Configuration.Language,
		Pushed:   c.Configuration.Pushed,
		Min:      c.Configuration.Min,
		Max:      c.Configuration.Max,
		GPool:    c.GPool,
	})
	srt := NewRetryTask(stc, scheduler.TaskFunc(SearchTaskFunc))
	c.SPool.Schedule(NewRetryTaskContext(&RetryTaskInfo{
		Times: 3,
	}), srt)
	c.SPool.Wait()
	logrus.Infof("Search task finished")

	type r struct {
		FullName string `json:"full_name"`
	}
	rs, err := func() ([]r, error) {
		rs := make([]r, 0)
		file, err := filetool.Open(c.StorePath.Repo+ReposJSON, filetool.RDONLY, 0644)
		if err != nil {
			return nil, err
		}
		return rs, filetool.NewDecoder(file).Decode(&rs)
	}()
	if err != nil {
		logrus.Errorf("Getting repos failed, %v", err)
		return err
	}

	infof, _ := sync.Open("")
	defer infof.Close()
	defer func() {
		err := infof.Save(c.StorePath.Repo + InfoJSON)
		if err != nil {
			logrus.Errorf("Save repositories' information failed, %v", err)
		}
	}()

	for index := range rs {
		i := strings.Index(rs[index].FullName, "/")
		if i == -1 || i+2 > len(rs[index].FullName) {
			logrus.Errorf("Invalid full name: %s", rs[index])
			continue
		}
		ltc := NewListTaskContext(&ListTaskInfo{
			Dir:   c.StorePath.Cache,
			URL:   "github.com/" + rs[index].FullName,
			Owner: rs[index].FullName[:i],
			Repo:  rs[index].FullName[i+1:],
			Info:  infof,
			GPool: c.GPool,
		})
		lrt := NewRetryTask(ltc, scheduler.TaskFunc(ListTaskFunc))
		c.SPool.Schedule(NewRetryTaskContext(&RetryTaskInfo{
			Times: 3,
		}), lrt)
	}
	c.SPool.Wait()
	logrus.Info("List task finished")

	itc := NewIndexTaskContext(&IndexTaskInfo{
		CacheDir: c.StorePath.Cache,
		ReposDir: c.StorePath.Repos,
		Info:     infof,
	})
	irt := NewRetryTask(itc, scheduler.TaskFunc(IndexTaskFunc))
	c.SPool.Schedule(NewRetryTaskContext(&RetryTaskInfo{
		Times: 3,
	}), irt)
	c.SPool.Wait()
	logrus.Info("Index task finished")

	return nil
}
