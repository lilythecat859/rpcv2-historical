// SPDX-License-Identifier: AGPL-3.0-only
package api

import (
	"crypto/tls"
	"net/http"
)

func (s *Server) Apply(opts ...Option) {
	for _, o := range opts {
		o(s)
	}
}

type Option func(*Server)

func TLS(cert, key string) Option {
	return func(s *Server) {
		if cert == "" || key == "" {
			return
		}
		cfg := &tls.Config{MinVersion: tls.VersionTLS13}
		srv := &http.Server{TLSConfig: cfg}
		srv.Handler = s
	}
}
