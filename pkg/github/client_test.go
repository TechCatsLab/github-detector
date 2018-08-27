/*
 * Revision History:
 *     Initial: 2018/07/18        Li Zebang
 */

package github

import (
	"testing"
)

func TestIsTokenValid_Normal(t *testing.T) {
	if err := IsTokenValid(&Token{Tag: "1", Token: "1"}); err != nil {
		t.Errorf("TestIsTokenValid_Normal() error = %v\n", err)
	}
}

func TestIsTokenValid_Error(t *testing.T) {
	if err := IsTokenValid(nil); err != ErrInvalidToken {
		t.Errorf("TestIsTokenValid_Error() error = %v\n", err)
	}
	if err := IsTokenValid(&Token{Tag: "1", Token: ""}); err != ErrInvalidToken {
		t.Errorf("TestIsTokenValid_Error() error = %v\n", err)
	}
}

func TestNewClient_Normal(t *testing.T) {
	if client, err := NewClient(&Token{Tag: "1", Token: "1"}); err != nil || client.Tag != "1" {
		t.Errorf("TestNewClient_Normal() error = %v\n", err)
	}
}

func TestNewClient_Error(t *testing.T) {
	if client, err := NewClient(nil); err != ErrInvalidToken || client != nil {
		t.Errorf("TestNewClient_Error() error = %v\n", err)
	}
	if client, err := NewClient(&Token{Tag: "1", Token: ""}); err != ErrInvalidToken || client != nil {
		t.Errorf("TestNewClient_Error() error = %v\n", err)
	}
}

func TestNewClient_DefualtTag(t *testing.T) {
	if client, err := NewClient(&Token{Tag: "", Token: "1"}); err != nil || client.Tag != DefualtClientTag {
		t.Errorf("TestNewClient_Error() error = %v\n", err)
	}
}
