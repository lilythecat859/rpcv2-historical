// SPDX-License-Identifier: AGPL-3.0-only
package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rpcv2-historical/internal/security"
)

type Server struct {
	db  ClickHouse
	key security.Signer
	val security.Validator
}

type ClickHouse interface {
	GetBlock(ctx context.Context, slot domain.Slot) (*domain.Block, error)
	GetTransaction(ctx context.Context, tx string) (*domain.Transaction, error)
	GetSigsForAddress(ctx context.Context, addr string, limit int) ([]domain.SigForAddr, error)
}

func NewServer(db ClickHouse, key security.Signer, val security.Validator) *Server {
	return &Server{db: db, key: key, val: val}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		http.Error(w, "only POST", http.StatusMethodNotAllowed)
		return
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, Response{JSONRPC: "2.0", ID: nil, Error: ErrParse})
		return
	}
	handler, ok := routes[req.Method]
	if !ok {
		writeJSON(w, Response{JSONRPC: "2.0", ID: req.ID, Error: ErrMethodNF})
		return
	}
	// ACL check via JWT scope
	scope, err := security.BearerScope(s.val, r)
	if err != nil {
		writeJSON(w, Response{JSONRPC: "2.0", ID: req.ID, Error: &Error{Code: -32003, Message: err.Error()}})
		return
	}
	ctx = context.WithValue(ctx, "scope", scope)
	result, err := handler(ctx, req.Params)
	var resp Response
	if err != nil {
		resp = Response{JSONRPC: "2.0", ID: req.ID, Error: &Error{Code: -32000, Message: err.Error()}}
	} else {
		resp = Response{JSONRPC: "2.0", ID: req.ID, Result: result}
	}
	writeJSON(w, resp)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
