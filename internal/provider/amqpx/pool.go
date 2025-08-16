//go:build !mock
// +build !mock

package amqpx

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrAlreadyClosed = errors.New("already closed")

type entry struct {
	conn     *Connection
	channels []*Channel
}

type Pool struct {
	url     string
	size    int
	index   int
	closed  int32
	rwMtx   sync.RWMutex
	entries map[int]*entry
	ctor    func() (*Connection, error)
}

func DialPool(url string, size int) *Pool {
	return &Pool{
		url:     url,
		size:    size,
		index:   -1,
		rwMtx:   sync.RWMutex{},
		entries: map[int]*entry{},
		ctor: func() (*Connection, error) {
			return Dial(url)
		},
	}
}

func (p *Pool) Channel() (*Channel, error) {
	defer p.rwMtx.Unlock()
	p.rwMtx.Lock()

	_, ok := p.entries[p.index]
	if !ok {
		// entry not found create new connection
		conn, err := p.ctor()
		if err != nil {
			return nil, err
		}

		p.index += 1
		p.entries[p.index] = &entry{conn: conn}
	}

	if len(p.entries[p.index].channels)+1 > p.size {
		// entry count reach limit
		conn, err := p.ctor()
		if err != nil {
			return nil, err
		}

		p.index += 1
		p.entries[p.index] = &entry{conn: conn}
	}

	channel, err := p.entries[p.index].conn.Channel()
	if err != nil {
		return nil, err
	}
	p.entries[p.index].channels = append(p.entries[p.index].channels, channel)

	return channel, nil
}

func (p *Pool) Close() error {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		// pool already closed
		return ErrAlreadyClosed
	}

	for _, e := range p.entries {
		for _, c := range e.channels {
			_ = c.Close()
		}
		_ = e.conn.Close()
	}

	return nil
}

func (p *Pool) ConnectionCount() int {
	defer p.rwMtx.RUnlock()
	p.rwMtx.RLock()

	return len(p.entries)
}

func (p *Pool) ChannelCount() int {
	defer p.rwMtx.RUnlock()
	p.rwMtx.RLock()

	var count int
	for _, item := range p.entries {
		count += len(item.channels)
	}

	return count
}
