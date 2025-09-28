// SPDX-License-Identifier: AGPL-3.0-only
package cache

import (
	"sync"
	"time"
)

type TTLShard struct {
	mu    sync.RWMutex
	items map[string]item
	ttl   time.Duration
}

type item struct {
	v      interface{}
	expire int64
}

func NewShard(ttl time.Duration) *TTLShard {
	s := &TTLShard{items: make(map[string]item), ttl: ttl}
	go s.reap()
	return s
}

func (s *TTLShard) Get(k string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[k]
	if !ok || time.Now().UnixNano() > it.expire {
		return nil, false
	}
	return it.v, true
}

func (s *TTLShard) Set(k string, v interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[k] = item{v: v, expire: time.Now().Add(s.ttl).UnixNano()}
}

func (s *TTLShard) reap() {
	tk := time.NewTicker(s.ttl)
	defer tk.Stop()
	for range tk.C {
		now := time.Now().UnixNano()
		s.mu.Lock()
		for k, v := range s.items {
			if now > v.expire {
				delete(s.items, k)
			}
		}
		s.mu.Unlock()
	}
}
