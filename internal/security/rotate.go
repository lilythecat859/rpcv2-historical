// SPDX-License-Identifier: AGPL-3.0-only
package security

import (
	"context"
	"log"
	"sync"
	"time"
)

type Rotator struct {
	key   *edKey
	Every time.Duration
	stop  chan struct{}
	wg    sync.WaitGroup
}

func NewRotator(k *edKey, every time.Duration) *Rotator {
	return &Rotator{key: k, Every: every}
}

func (r *Rotator) Start() {
	r.stop = make(chan struct{})
	r.wg.Add(1)
	go r.run()
}

func (r *Rotator) Stop() {
	close(r.stop)
	r.wg.Wait()
}

func (r *Rotator) run() {
	defer r.wg.Done()
	tk := time.NewTicker(r.Every)
	defer tk.Stop()
	for {
		select {
		case <-tk.C:
			if err := r.key.Rotate(); err != nil {
				log.Printf("key rotation failed: %v", err)
				continue
			}
			log.Println("jwt key rotated")
		case <-r.stop:
			return
		}
	}
}
