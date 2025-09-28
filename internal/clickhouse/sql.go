// SPDX-License-Identifier: AGPL-3.0-only
package clickhouse

import (
	"context"
	"fmt"

	"github.com/rpcv2-historical/internal/domain"
)

func (db *DB) GetBlock(ctx context.Context, slot domain.Slot) (*domain.Block, error) {
	row := db.conn.QueryRowContext(ctx, `SELECT slot, blockhash, parent_slot, block_time, tx_count, updated_on FROM blocks WHERE slot=?`, slot)
	var b domain.Block
	err := row.Scan(&b.Slot, &b.Hash, &b.ParentSlot, &b.BlockTime, &b.TxCount, &b.UpdatedOn)
	if err != nil {
		return nil, fmt.Errorf("select block: %w", err)
	}
	return &b, nil
}

func (db *DB) GetTransaction(ctx context.Context, tx string) (*domain.Transaction, error) {
	row := db.conn.QueryRowContext(ctx, `SELECT slot, tx_hash, idx, meta, message, block_time, updated_on FROM transactions WHERE tx_hash=?`, tx)
	var t domain.Transaction
	err := row.Scan(&t.Slot, &t.TxHash, &t.Index, &t.Meta, &t.Message, &t.BlockTime, &t.UpdatedOn)
	if err != nil {
		return nil, fmt.Errorf("select tx: %w", err)
	}
	return &t, nil
}

func (db *DB) GetSigsForAddress(ctx context.Context, addr string, limit int) ([]domain.SigForAddr, error) {
	rows, err := db.conn.QueryContext(ctx, `SELECT address, signature, slot, memo, error, block_time, updated_on FROM sigs_for_address WHERE address=? ORDER BY slot DESC LIMIT ?`, addr, limit)
	if err != nil {
		return nil, fmt.Errorf("select sigs: %w", err)
	}
	defer rows.Close()
	var out []domain.SigForAddr
	for rows.Next() {
		var s domain.SigForAddr
		if err := rows.Scan(&s.Address, &s.TxHash, &s.Slot, &s.Memo, &s.Err, &s.BlockTime, &s.UpdatedOn); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}
