/*
 * Revision History:
 *     Initial: 2018/08/03        Li Zebang
 */

package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

// Search searches repositories via various criteria.
//
// https://developer.github.com/v3/search/#search-repositories
// https://help.github.com/articles/searching-repositories
func (c *Client) Search(language string, pushed time.Duration, min, max, page int) (*github.RepositoriesSearchResult, *github.Response, error) {
	c.bucket.Wait(1)

	return c.Client.Search.Repositories(context.Background(),
		fmt.Sprintf("language:%s+stars:%d..%d+pushed:>=%s&sort=stars&order=desc&page=%d&per_page=100",
			language, min, max, time.Now().Add(-pushed).Format("2006-01-02"), page), nil)
}

// List returns the contents of a file or directory in a repository.
//
// https://developer.github.com/v3/repos/contents/#get-contents
func (c *Client) List(owner, repo, path string) ([]*github.RepositoryContent, *github.Response, error) {
	c.bucket.Wait(1)

	_, dc, resp, err := c.Client.Repositories.GetContents(context.Background(), owner, repo, path, nil)

	return dc, resp, err
}
