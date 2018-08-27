/*
 * Revision History:
 *     Initial: 2018/07/18        Li Zebang
 */

package github

import (
	"context"
	"errors"

	"github.com/google/go-github/github"
	"github.com/juju/ratelimit"
	"golang.org/x/oauth2"
)

const (
	rate     = 5000.0 / 3600.0
	capacity = 5000

	// DefualtClientTag is used when the tag is "".
	DefualtClientTag = "defualt"
)

// Token contains tag and github.
type Token struct {
	Tag   string `json:"tag,omitempty"`
	Token string `json:"token"`
}

// Client contains a available GitHub Client.
type Client struct {
	*github.Client

	Tag    string
	bucket *ratelimit.Bucket
}

var (
	// ErrInvalidToken occurs when the token object is nil or token is "".
	ErrInvalidToken = errors.New("invalid token")
)

// IsTokenValid returns ErrInvalidToken if token is invalid.
func IsTokenValid(token *Token) error {
	if token == nil || token.Token == "" {
		return ErrInvalidToken
	}
	return nil
}

// NewClient returns a Client.
func NewClient(token *Token) (*Client, error) {
	err := IsTokenValid(token)
	if err != nil {
		return nil, err
	}
	if token.Tag == "" {
		token.Tag = DefualtClientTag
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.Token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := &Client{
		Client: github.NewClient(tc),
		Tag:    token.Tag,
		bucket: ratelimit.NewBucketWithRate(rate, capacity),
	}
	client.bucket.TakeAvailable(capacity)

	return client, nil
}
