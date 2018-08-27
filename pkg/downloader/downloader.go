/*
 * Revision History:
 *     Initial: 2018/08/02        Li Zebang
 */

package downloader

import (
	"io"
	"net"
	"net/http"
	"time"
)

// A Client is an HTTP client.
type Client struct {
	client *http.Client
}

func transport() http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// DefaultTimeout is used if the parameter timeout less then zero.
const DefaultTimeout = 30 * time.Second

// NewClient returns an HTTP client.
func NewClient(timeout time.Duration) *Client {
	if timeout < 0 {
		timeout = DefaultTimeout
	}

	return &Client{client: &http.Client{Transport: transport(), Timeout: timeout}}
}

// Download gets the vendor profile and saves it to target file.
func (c *Client) Download(url string, buf io.Writer) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	_, err = io.Copy(buf, resp.Body)
	return err
}

// Get -
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}
