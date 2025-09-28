// SPDX-License-Identifier: AGPL-3.0-only
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rpcv2-historical/internal/clickhouse"
)

var (
	dbDSN   = flag.String("db", "clickhouse://default:@localhost:9000/rpcv2_hist", "")
	rows    = flag.Uint64("rows", 2_300_000_000, "rows to generate")
	workers = flag.Int("workers", 16, "parallel workers")
	qps     = flag.Int("qps", 20_000, "target queries/sec")
)

func main() {
	flag.Parse()
	db, err := clickhouse.New(*dbDSN)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	ctx := context.Background()
	fmt.Println("seeding...")
	seed(ctx, db)
	fmt.Println("benchmarking...")
	bench(ctx, db)
}

func seed(ctx context.Context, db *clickhouse.DB) {
	var wg sync.WaitGroup
	wg.Add(*workers)
	rowsPerWorker := *rows / uint64(*workers)
	for w := 0; w < *workers; w++ {
		go func(id int) {
			defer wg.Done()
			for i := uint64(0); i < rowsPerWorker; i++ {
				addr := fmt.Sprintf("Addr%d%d", id, i)
				sig := fmt.Sprintf("Sig%d%d%d", id, i, rand.Uint64())
				if err := db.InsertSig(ctx, addr, sig); err != nil {
					log.Printf("insert: %v", err)
				}
			}
		}(w)
	}
	wg.Wait()
}

func bench(ctx context.Context, db *clickhouse.DB) {
	var ops uint64
	interval := 10 * time.Second
	start := time.Now()
	tick := time.NewTicker(interval)
	defer tick.Stop()
	go func() {
		for range tick.C {
			elapsed := time.Since(start).Seconds()
			o := atomic.SwapUint64(&ops, 0)
			fmt.Printf("QPS: %.0f\n", float64(o)/elapsed)
		}
	}()
	// query loop
	for {
		addr := fmt.Sprintf("Addr%d", rand.Intn(*workers))
		_, _ = db.GetSigsForAddress(ctx, addr, 1000)
		atomic.AddUint64(&ops, 1)
	}
}
