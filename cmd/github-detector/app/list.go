/*
 * Revision History:
 *     Initial: 2018/08/04        Li Zebang
 */

package app

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/TechCatsLab/logging/logrus"
	github "github.com/google/go-github/github"

	"github.com/fengyfei/github-detector/pkg/downloader"
	"github.com/fengyfei/github-detector/pkg/filetool"
	pool "github.com/fengyfei/github-detector/pkg/github"
	store "github.com/fengyfei/github-detector/pkg/store/file/yaml"
)

type (
	// ListTaskInfo -
	ListTaskInfo struct {
		Dir   string
		URL   string
		Owner string
		Repo  string
		Path  string

		Info *store.File

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
		rc, err := client.List(info.Owner, info.Repo, info.Path)
		if err != nil {
			logrus.Errorf("LIST %s FAILED, ERROR -- %v", info.URL, err)
			return err
		}
		dt := filter(rc, info, repo)

		if dt.Type == NoVendorType {
			repo.Type = "no vendor"
			repo.Status = "FINISH"
			logrus.Infof("FINISH: %s, NO VENDOR", info.URL)
			return nil
		}

		if dt.Type == NotSupportType {
			repo.Type = "not support"
			repo.Status = "FINISH"
			logrus.Infof("FINISH: %s, NOT SUPPORT", info.URL)
			return nil
		}

		if dt.Type == VendorType {
			src, err := download(dt)
			if err != nil {
				repo.Status = "FAIL: " + err.Error()
				logrus.Errorf("FAIL: %s, ERROR -- %v", info.URL, err)
				return err
			}
			src.URL = info.URL
			src.Type = repo.Type

			err = cache(info.Dir+strings.Replace(info.URL, "/", "-", -1), src)
			if err != nil {
				logrus.Errorf("CACHE FAILED -- %s", info.URL)
				return err
			}
			logrus.Infof("CACHE SUCCESS -- %s", info.URL)

			repo.Status = "SUCCESS"
			logrus.Infof("SUCCESS: %s", info.URL)
			return nil
		}
	}
}

func filter(rc []*github.RepositoryContent, info *ListTaskInfo, repo *Info) *DownloadTask {
	dt := &DownloadTask{}
	for index := range rc {
		switch rc[index].GetName() {
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
			info.Path = "Godeps"
			repo.Type = "godep"
			dt.Type = AgainType
			return dt
		case "Godeps.json":
			dt.Lock = rc[index].GetDownloadURL()
			dt.Type = VendorType
			repo.Type = "godep"
		case "vendor":
			if dt.Type == IgnoreType {
				info.Path = "vendor"
				dt.Type = AgainType
			}
		case "vendor.json":
			dt.Lock = rc[index].GetDownloadURL()
			dt.Type = VendorType
			repo.Type = "govendor"
		}
	}
	if dt.Type == IgnoreType && info.Path != "vendor" {
		dt.Type = NoVendorType
	}
	if dt.Type == IgnoreType && info.Path == "vendor" {
		dt.Type = NotSupportType
	}
	return dt
}

func download(dt *DownloadTask) (*Source, error) {
	src := &Source{}
	client := downloader.NewClient(downloader.DefaultTimeout)

	if dt.Type != VendorType {
		return nil, nil
	}

	if dt.Conf != "" {
		resp, err := client.Get(dt.Conf)
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
		resp, err := client.Get(dt.Lock)
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
