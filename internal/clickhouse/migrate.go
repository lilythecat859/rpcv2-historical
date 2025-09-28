// SPDX-License-Identifier: AGPL-3.0-only
package clickhouse

import (
	"context"
	"embed"
	"fmt"
	"strings"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func (db *DB) Migrate(ctx context.Context) error {
	files := []string{
		"migrations/001_init_blocks.sql",
		"migrations/002_init_transactions.sql",
		"migrations/003_init_sigs_for_address.sql",
	}
	conn, err := db.conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	for _, name := range files {
		b, err := migrationsFS.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		for _, stmt := range strings.Split(string(b), ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if err := conn.Exec(ctx, stmt); err != nil {
				return fmt.Errorf("exec %s: %w", name, err)
			}
		}
	}
	return nil
}
