/*
 * Revision History:
 *     Initial: 2018/08/04        Li Zebang
 */

package app

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	github "github.com/google/go-github/github"

	"github.com/TechCatsLab/logging/logrus"

	"github.com/TechCatsLab/github-detector/pkg/filetool"
	pool "github.com/TechCatsLab/github-detector/pkg/github"
	"github.com/TechCatsLab/github-detector/pkg/sync"
)

type (
	// ListTaskInfo -
	ListTaskInfo struct {
		Dir   string
		URL   string
		Owner string
		Repo  string
		Path  string

		Info *sync.Map

		GPool pool.Pool
	}

	// Info -
	Info struct {
		URL    string `yaml:"url"`
		Type   string `yaml:"type"`
		Status string `yaml:"status"`
	}

	// DownloadType -
	DownloadType int8

	// DownloadTask -
	DownloadTask struct {
		Type DownloadType
		Conf string
		Lock string
	}

	// Source -
	Source struct {
		URL  string `json:"url"`
		Type string `json:"type"`
		Conf []byte `json:"conf,omitempty"`
		Lock []byte `json:"lock"`
	}
)

const (
	// IgnoreType -
	IgnoreType DownloadType = iota
	// NoVendorType -
	NoVendorType
	// NotSupportType -
	NotSupportType
	// AgainType -
	AgainType
	// VendorType -
	VendorType
)

// NewListTaskContext -
func NewListTaskContext(info *ListTaskInfo) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ListTaskKey, info)

	return ctx
}

// ListTaskFunc -
func ListTaskFunc(ctx context.Context) error {
	info, ok := ctx.Value(ListTaskKey).(*ListTaskInfo)
	if !ok {
		return errors.New("assertion fail")
	}

	client := info.GPool.Get(pool.DefualtClientTag)
	defer info.GPool.Put(client)
	if client == nil {
		return errors.New("no available client")
	}

	repo := &Info{
		URL: info.URL,
	}

	defer func() {
		info.Info.Upsert(repo.URL, repo)
	}()

	for {
		rc, resp, err := client.List(info.Owner, info.Repo, info.Path)
		if err != nil && resp.Remaining == 0 {
			logrus.Errorf("List %s failed, error -- %v", info.URL, err)
			<-time.After(time.Until(resp.Reset.Time))
			return err
		} else if err != nil {
			logrus.Errorf("List %s failed, error -- %v", info.URL, err)
			return err
		}
		dt := filter(rc, info, repo)

		if dt.Type == NoVendorType {
			repo.Type = "no vendor"
			repo.Status = "Finish"
			logrus.Infof("Finish: %s, no vendor", info.URL)
			return nil
		}

		if dt.Type == NotSupportType {
			repo.Type = "not support"
			repo.Status = "Finish"
			logrus.Infof("Finish: %s, not support", info.URL)
			return nil
		}

		if dt.Type == VendorType {
			src, err := download(dt)
			if err != nil {
				repo.Status = "Fail: " + err.Error()
				logrus.Errorf("Fail: %s, error -- %v", info.URL, err)
				return err
			}
			src.URL = info.URL
			src.Type = repo.Type

			err = cache(info.Dir+strings.Replace(info.URL, "/", "-", -1), src)
			if err != nil {
				logrus.Errorf("Cache failed -- %s", info.URL)
				return err
			}
			logrus.Infof("Cache success -- %s", info.URL)

			repo.Status = "SUCCESS"
			logrus.Infof("Success: %s", info.URL)
			return nil
		}
	}
}

func filter(rc []*github.RepositoryContent, info *ListTaskInfo, repo *Info) *DownloadTask {
	dt := &DownloadTask{}
	for index := range rc {
		switch rc[index].GetPath() {
		case "Gopkg.lock":
			dt.Lock = rc[index].GetDownloadURL()
			dt.Type = VendorType
			repo.Type = "dep"
		case "Gopkg.toml":
			dt.Conf = rc[index].GetDownloadURL()
			dt.Type = VendorType
			repo.Type = "dep"
		case "glide.lock":
			dt.Lock = rc[index].GetDownloadURL()
			dt.Type = VendorType
			repo.Type = "glide"
		case "glide.yaml":
			dt.Conf = rc[index].GetDownloadURL()
			dt.Type = VendorType
			repo.Type = "glide"
		case "Godeps":
			info.Path = rc[index].GetPath()
			repo.Type = "godep"
			dt.Type = AgainType
			return dt
		case "Godeps/Godeps.json":
			dt.Lock = rc[index].GetDownloadURL()
			dt.Type = VendorType
			repo.Type = "godep"
		case "vendor":
			if dt.Type == IgnoreType {
				info.Path = rc[index].GetPath()
				dt.Type = AgainType
			}
		case "vendor/vendor.json":
			dt.Lock = rc[index].GetDownloadURL()
			dt.Type = VendorType
			repo.Type = "govendor"
		}
	}
	if dt.Type == IgnoreType && info.Path == "vendor" {
		dt.Type = NotSupportType
	}
	if dt.Type == IgnoreType && info.Path != "vendor" {
		dt.Type = NoVendorType
	}
	return dt
}

func download(dt *DownloadTask) (*Source, error) {
	src := &Source{}

	if dt.Type != VendorType {
		return nil, nil
	}

	if dt.Conf != "" {
		resp, err := http.Get(dt.Conf)
		if err != nil {
			return nil, err
		}
		src.Conf, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
	}

	if dt.Lock != "" {
		resp, err := http.Get(dt.Lock)
		if err != nil {
			return nil, err
		}
		src.Lock, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
	}

	return src, nil
}

func cache(path string, src *Source) error {
	f, err := filetool.Open(path, filetool.TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return filetool.NewEncoder(f).Encode(src)
}
