/*
 * Revision History:
 *     Initial: 2018/07/18        Li Zebang
 */

package github

import (
	"testing"
)

func Test_circle(t *testing.T) {
	clients := []*Client{&Client{}, &Client{}, &Client{}, &Client{}, &Client{}}
	circle := circle{clients: make([]*Client, len(clients))}
	for index := 0; index < len(clients); index++ {
		circle.put(clients[index])
	}

	for index := 0; index < len(clients)-1; index++ {
		if circle.clients[index] != clients[index] {
			t.Errorf("Test_circle() circle.clients[%d] = %v, not %v\n", index, circle.clients[index], clients[index])
		}
	}
	if circle.clients[len(clients)-1] != nil {
		t.Errorf("Test_circle() circle.clients[%d] = %v, not nil\n", len(clients)-1, circle.clients[len(clients)-1])
	}

	for index := 0; index < len(clients)-1; index++ {
		if circle.get() != clients[index] {
			t.Errorf("Test_circle() circle.clients[%d] = %v, not %v\n", index, circle.clients[index], clients[index])
		}
	}
	if circle.get() != nil {
		t.Errorf("Test_circle() circle.clients[%d] = %v, not nil\n", len(clients)-1, circle.clients[len(clients)-1])
	}
}

func TestNewPool_Normal(t *testing.T) {
	if pool := NewPool(&Token{Tag: "1", Token: "1"}, &Token{Tag: "1", Token: "1"}, &Token{Tag: "1", Token: "1"}); pool == nil {
		t.Errorf("NewPool() = %v\n", pool)
	}

	if pool := NewPool(&Token{Token: "1"}, &Token{Token: "1"}, &Token{Token: "1"}); pool == nil {
		t.Errorf("NewPool() = %v\n", pool)
	}
}

func TestNewPool_WantErr(t *testing.T) {
	if pool := NewPool(); pool != nil {
		t.Errorf("NewPool() = %v\n", pool)
	}

	if pool := NewPool(&Token{Tag: "1", Token: ""}, &Token{Tag: "1", Token: ""}, &Token{Tag: "1", Token: ""}); pool != nil {
		t.Errorf("NewPool() = %v\n", pool)
	}
}

func Test_pool_Get_Put(t *testing.T) {
	pool := NewPool(&Token{Token: "1"}, &Token{Token: "1"}, &Token{Token: "1"}, &Token{Token: "1"}, &Token{Token: "1"})
	clients := make([]*Client, 5)
	for index := range clients {
		clients[index] = pool.Get("defualt")
	}
	for index := range clients {
		pool.Put(clients[index])
	}

	if client := pool.Get("0"); client != nil {
		t.Errorf("Test_pool_Get_Put() pool.Get = %v, not nil\n", client)
	}

	for index := range clients {
		if client := pool.Get("defualt"); client != clients[index] {
			t.Errorf("Test_pool_Get_Put() pool.Get = %v, not client[%d] = %v\n", client, index, clients[index])
		}
	}

	if client := pool.Get("defualt"); client != nil {
		t.Errorf("Test_pool_Get_Put() pool.Get = %v, not nil\n", client)
	}
}
