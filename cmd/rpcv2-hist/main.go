// SPDX-License-Identifier: AGPL-3.0-only
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rpcv2-historical/internal/api"
	"github.com/rpcv2-historical/internal/clickhouse"
	"github.com/rpcv2-historical/internal/security"
)

var (
	listen   = flag.String("listen", ":8899", "HTTP listen address")
	dbDSN    = flag.String("db", "clickhouse://default:@localhost:9000/rpcv2_hist", "DB DSN")
	jwtKey   = flag.String("jwt-key", "", "Ed25519 private key (base64) for signing JWTs")
	certFile = flag.String("cert", "", "TLS cert (leave empty for auto-self-signed)")
	keyFile  = flag.String("key", "", "TLS key")
)

func main() {
	flag.Parse()
	if *jwtKey == "" {
		log.Println("=== RPCv2 Historical quick-start ===")
		log.Println("1. Generate a keypair:")
		log.Println("   ./scripts/gen-keypair.sh")
		log.Println("2. Re-run with -jwt-key=<base64>")
		return
	}
	key, err := security.ParseEdKey(*jwtKey)
	if err != nil {
		log.Fatalf("bad jwt-key: %v", err)
	}
	db, err := clickhouse.New(*dbDSN)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	if err := db.Migrate(context.Background()); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	// optional key rotator
	rot := security.NewRotator(key, 24*time.Hour)
	rot.Start()
	defer rot.Stop()

	srv := api.NewServer(db, key, key)
	srv.Apply(api.TLS(*certFile, *keyFile))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("listening on %s", *listen)
		if err := srv.Start(*listen); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http: %v", err)
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
