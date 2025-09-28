// SPDX-License-Identifier: AGPL-3.0-only
package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rpcv2-historical/internal/domain"
)

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrParse       = &Error{-32700, "parse error"}
	ErrInvalidReq  = &Error{-32600, "invalid request"}
	ErrMethodNF    = &Error{-32601, "method not found"}
	ErrInvalidParams = &Error{-32602, "invalid params"}
	ErrInternal    = &Error{-32603, "internal error"}
)

type Handler func(ctx context.Context, params json.RawMessage) (interface{}, error)

var routes = map[string]Handler{
	"getBlock":              handleGetBlock,
	"getTransaction":        handleGetTransaction,
	"getSignaturesForAddress": handleSigsForAddr,
	"getBlocksWithLimit":    handleBlocksWithLimit,
	"getBlockTime":          handleGetBlockTime,
	"getSlot":               handleGetSlot,
}

func handleGetBlock(ctx context.Context, p json.RawMessage) (interface{}, error) {
	var params []interface{}
	if err := json.Unmarshal(p, &params); err != nil || len(params) < 1 {
		return nil, ErrInvalidParams
	}
	slot := domain.Slot(params[0].(float64))
	return map[string]interface{}{"slot": slot}, nil
}

func handleGetTransaction(ctx context.Context, p json.RawMessage) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func handleSigsForAddr(ctx context.Context, p json.RawMessage) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func handleBlocksWithLimit(ctx context.Context, p json.RawMessage) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func handleGetBlockTime(ctx context.Context, p json.RawMessage) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func handleGetSlot(ctx context.Context, p json.RawMessage) (interface{}, error) {
	return nil, errors.New("not implemented")
}
