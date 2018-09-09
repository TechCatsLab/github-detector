/*
 * Revision History:
 *     Initial: 2018/08/26        Li Zebang
 */

package app

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/TechCatsLab/logging/logrus"

	"github.com/TechCatsLab/github-detector/pkg/codec"
	"github.com/TechCatsLab/github-detector/pkg/codec/conf"
	"github.com/TechCatsLab/github-detector/pkg/codec/lock"
	"github.com/TechCatsLab/github-detector/pkg/filetool"
	"github.com/TechCatsLab/github-detector/pkg/sync"
)

type (
	// IndexTaskInfo -
	IndexTaskInfo struct {
		CacheDir string
		ReposDir string

		Info *sync.Map
	}

	// Dep -
	Dep struct {
		URL      string    `json:"url"`
		Type     string    `json:"type"`
		Conf     conf.File `json:"conf"`
		Lock     lock.File `json:"lock"`
		Direct   []string  `json:"direct"`
		Indirect []string  `json:"indirect"`
	}
)

// NewIndexTaskContext -
func NewIndexTaskContext(info *IndexTaskInfo) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, IndexTaskKey, info)

	return ctx
}

// IndexTaskFunc -
func IndexTaskFunc(ctx context.Context) error {
	info, ok := ctx.Value(IndexTaskKey).(*IndexTaskInfo)
	if !ok {
		return errors.New("assertion fail")
	}

	files, err := ioutil.ReadDir(info.CacheDir)
	if err != nil {
		logrus.Errorf("Read cache directory failed, %v", err)
	}

	deps := make(map[string]*Dep, len(files))
	defer func() {
		for key, val := range deps {
			file, err := filetool.Open(info.ReposDir+strings.Replace(key, "/", "-", -1), filetool.TRUNC, 0644)
			if err != nil {
				logrus.Errorf("Store %s failed, %v", key, err)
				info.Info.Upsert(key, &Info{URL: val.URL, Type: val.Type, Status: "FAIL: " + err.Error()})
				continue
			}
			err = filetool.NewEncoder(file).Encode(val)
			file.Close()
			if err != nil {
				logrus.Errorf("Store %s failed, %v", key, err)
				info.Info.Upsert(key, &Info{URL: val.URL, Type: val.Type, Status: "FAIL: " + err.Error()})
				continue
			}
		}
	}()

	for index := range files {
		name := files[index].Name()
		file, err := filetool.Open(info.CacheDir+name, filetool.RDONLY, 0644)
		if err != nil {
			logrus.Errorf("Open %s failed, %v", err)
			continue
		}

		var src Source
		err = filetool.NewDecoder(file).Decode(&src)
		file.Close()
		if err != nil {
			logrus.Errorf("Decode %s failed, %v", err)
			continue
		}

		var c codec.Codec
		switch src.Type {
		case "dep":
			c = codec.Dep()
		case "glide":
			c = codec.Gilde()
		case "godep":
			c = codec.Godep()
		case "govendor":
			c = codec.Govendor()
		}

		dep, ok := deps[src.URL]
		if ok {
			dep.URL = src.URL
			dep.Type = src.Type
		} else {
			dep = &Dep{URL: src.URL, Type: src.Type}
			deps[src.URL] = dep
		}
		if len(src.Conf) != 0 {
			dep.Conf, err = c.ParseConfFile(src.Conf)
			if err != nil {
				logrus.Errorf("Parse conf file %s failed, %v", src.URL, err)
				info.Info.Upsert(src.URL, &Info{URL: src.URL, Type: src.Type, Status: "FAIL: " + err.Error()})
			} else {
				direct := dep.Conf.Deps()
				for _, value := range direct {
					d, ok := deps[value]
					if ok {
						d.Direct = append(d.Direct, src.URL)
					} else {
						d = &Dep{URL: value, Direct: []string{src.URL}}
					}
					deps[value] = d
				}
			}
		}
		if len(src.Lock) != 0 {
			dep.Lock, err = c.ParseLockFile(src.Lock)
			if err != nil && err != io.EOF {
				logrus.Errorf("Parse lock file %s failed, %v", src.URL, err)
				info.Info.Upsert(src.URL, &Info{URL: src.URL, Type: src.Type, Status: "FAIL: " + err.Error()})
			} else {
				indirect := dep.Lock.Repos()
				for _, value := range indirect {
					d, ok := deps[value]
					if ok {
						d.Indirect = append(d.Indirect, src.URL)
					} else {
						d = &Dep{URL: value, Indirect: []string{src.URL}}
					}
					deps[value] = d
				}
			}
		}
	}

	return nil
}
