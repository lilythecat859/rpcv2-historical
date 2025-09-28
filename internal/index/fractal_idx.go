// SPDX-License-Identifier: AGPL-3.0-only
package index

import (
	"hash/fnv"
)

type FractalIdx struct{ shifts uint8 }

func New() *FractalIdx { return &FractalIdx{shifts: 12} }

func (f *FractalIdx) Shard(slot uint64) uint32 {
	return uint32(slot >> f.shifts)
}

func (f *FractalIdx) Partition(addr string, shards uint32) uint32 {
	h := fnv.New32a()
	h.Write([]byte(addr))
	return h.Sum32() % shards
}
