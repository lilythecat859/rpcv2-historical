// SPDX-License-Identifier: AGPL-3.0-only
package clickhouse

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type DB struct {
	conn *sql.DB
}

func New(dsn string) (*DB, error) {
	c, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := c.PingContext(ctx); err != nil {
		return nil, err
	}
	return &DB{conn: c}, nil
}

func (db *DB) Close() error { return db.conn.Close() }
