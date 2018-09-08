/*
 * Revision History:
 *     Initial: 2018/08/04        Li Zebang
 */

package app

import (
	"context"
	"errors"
	"time"

	github "github.com/google/go-github/github"

	"github.com/TechCatsLab/logging/logrus"

	"github.com/TechCatsLab/github-detector/pkg/filetool"
	pool "github.com/TechCatsLab/github-detector/pkg/github"
)

// SearchTaskInfo -
type SearchTaskInfo struct {
	RepoDir  string
	Language string
	Pushed   time.Duration
	Min      int
	Max      int

	GPool pool.Pool
}

// NewSearchTaskContext -
func NewSearchTaskContext(info *SearchTaskInfo) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, SearchTaskKey, info)

	return ctx
}

// SearchTaskFunc -
func SearchTaskFunc(ctx context.Context) (err error) {
	info, ok := ctx.Value(SearchTaskKey).(*SearchTaskInfo)
	if !ok {
		return errors.New("assert fail")
	}

	client := info.GPool.Get(pool.DefualtClientTag)
	defer info.GPool.Put(client)
	if client == nil {
		return errors.New("no available client")
	}

	// store stores
	store := func(dir string, repos *[]github.Repository) error {
		f, err := filetool.Open(dir, filetool.TRUNC, 0)
		if err != nil {
			return err
		}
		defer f.Close()

		err = filetool.NewEncoder(f).Encode(*repos)
		if err != nil {
			return err
		}
		return nil
	}

	// specify gets right index.
	specify := func(pre []github.Repository, rsr *github.RepositoriesSearchResult) int {
		index := 0
		for len(pre) != 0 && index < len(rsr.Repositories) {
			switch {
			case pre[len(pre)-1].GetFullName() == rsr.Repositories[index].GetFullName():
				index++
				return index
			case pre[len(pre)-1].GetStargazersCount() == rsr.Repositories[index].GetStargazersCount():
				index++
				break
			default:
				return index
			}
		}
		return index
	}

	rsr := &github.RepositoriesSearchResult{}
	repos := []github.Repository(nil)
	defer func() {
		if err == nil {
			err = store(info.RepoDir+ReposJSON, &repos)
		}
	}()

	for {
		rsr, err = client.Search(info.Language, info.Pushed, info.Min, info.Max, 1)
		if err != nil {
			return err
		} else if rsr.GetIncompleteResults() {
			err = errors.New("search results incomplete")
			return
		}
		index := specify(repos, rsr)
		repos = append(repos, rsr.Repositories[index:]...)
		logrus.Infof("Search number: %d", len(repos))
		if rsr.GetTotal() > 100 {
			info.Max = rsr.Repositories[len(rsr.Repositories)-1].GetStargazersCount()
			continue
		}
		return nil
	}
}
