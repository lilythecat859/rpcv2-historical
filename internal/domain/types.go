// SPDX-License-Identifier: AGPL-3.0-only
package domain

import "time"

type Slot uint64

type Transaction struct {
	Slot      Slot      `json:"slot"`
	TxHash    string    `json:"transaction_hash"`
	Index     uint32    `json:"index"`
	Meta      []byte    `json:"meta"`        // bincode blob
	Message   []byte    `json:"message"`     // bincode blob
	BlockTime int64     `json:"block_time"`
	UpdatedOn time.Time `json:"updated_on"`
}

type Block struct {
	Slot       Slot      `json:"slot"`
	Hash       string    `json:"blockhash"`
	ParentSlot Slot      `json:"parent_slot"`
	BlockTime  int64     `json:"block_time"`
	TxCount    uint32    `json:"transaction_count"`
	UpdatedOn  time.Time `json:"updated_on"`
}

type SigForAddr struct {
	Address   string    `json:"address"`
	TxHash    string    `json:"signature"`
	Slot      Slot      `json:"slot"`
	Memo      []byte    `json:"memo"`
	Err       *string   `json:"error,omitempty"`
	BlockTime int64     `json:"block_time"`
	UpdatedOn time.Time `json:"updated_on"`
}

// Bit-mask scopes
const (
	ScopeRead  uint32 = 1 << iota
	ScopeWrite
	ScopeAdmin
)
