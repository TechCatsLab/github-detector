/*
 * Revision History:
 *     Initial: 2018/08/03        Li Zebang
 */
package github

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/fengyfei/github-detector/pkg/filetool"
)

func TestClient_Search(t *testing.T) {
	client, err := NewClient(&Token{Token: "3c8bf037bd6bfe3858ad8649099b0ee82f7d6cf0"})
	if err != nil {
		t.Errorf("TestClient_Search() NewClient error: %s\n", err)
	}

	rsr, err := client.Search("golang", 365*24*time.Hour, 200, 400, 1)
	if err != nil {
		t.Errorf("TestClient_Search() client.Search error: %s\n", err)
	}

	f, err := filetool.Open("test.txt", filetool.TRUNC, 0)
	if err != nil {
		t.Errorf("TestClient_Search() file.Open error: %s\n", err)
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	i := interface{}(rsr.Repositories)
	err = encoder.Encode(i)
	if err != nil {
		t.Errorf("TestClient_Search() Encoder.Encode error: %s\n", err)
	}
}
