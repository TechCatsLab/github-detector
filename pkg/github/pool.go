/*
 * Revision History:
 *     Initial: 2018/07/18        Li Zebang
 */

package github

import (
	"log"
	"sync"
)

// Pool provides Get and Put function.
type Pool interface {
	Get(string) *Client
	Put(*Client)
}

type pool struct {
	locker  *sync.Mutex
	circles map[string]*circle
}

type circle struct {
	clients []*Client
	head    int
	rear    int
}

func (c *circle) get() *Client {
	if c.head == c.rear {
		return nil
	}

	client := c.clients[c.head]
	c.clients[c.head] = nil
	c.head = (c.head + 1) % len(c.clients)

	return client
}

func (c *circle) put(client *Client) {
	if (c.rear+1)%len(c.clients) == c.head {
		return
	}

	c.clients[c.rear] = client
	c.rear = (c.rear + 1) % len(c.clients)
}

// NewPool returns a GitHub Client Pool. It returns nil if there is no available Token.
func NewPool(tokens ...*Token) Pool {
	if len(tokens) == 0 {
		return nil
	}

	ts := make(map[string][]*Token, 0)
	for index := range tokens {
		err := IsTokenValid(tokens[index])
		if err != nil {
			log.Printf("%s: index %d -- %v", err, index, tokens[index])
			continue
		}
		if tokens[index].Tag == "" {
			tokens[index].Tag = DefualtClientTag
		}
		if _, exist := ts[tokens[index].Tag]; !exist {
			ts[tokens[index].Tag] = make([]*Token, 0, 1)
			ts[tokens[index].Tag] = append(ts[tokens[index].Tag], tokens[index])
		} else {
			ts[tokens[index].Tag] = append(ts[tokens[index].Tag], tokens[index])
		}
	}
	if len(ts) == 0 {
		return nil
	}

	p := &pool{
		locker:  &sync.Mutex{},
		circles: make(map[string]*circle, len(ts)),
	}
	for k, t := range ts {
		p.circles[k] = &circle{clients: make([]*Client, len(t)+1)}
		for _, token := range t {
			client, _ := NewClient(token)
			p.circles[token.Tag].put(client)
		}
	}

	return p
}

// Get a Client from Pool, it will return nil if there is no available Client.
func (p *pool) Get(tag string) *Client {
	p.locker.Lock()
	defer p.locker.Unlock()

	if tag == "" {
		tag = DefualtClientTag
	}

	if _, exist := p.circles[tag]; !exist {
		return nil
	}

	return p.circles[tag].get()
}

// Put the Client back to the Pool.
func (p *pool) Put(client *Client) {
	p.locker.Lock()
	defer p.locker.Unlock()

	if _, exist := p.circles[client.Tag]; exist {
		p.circles[client.Tag].put(client)
	}
}
