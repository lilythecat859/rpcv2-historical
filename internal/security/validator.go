// SPDX-License-Identifier: AGPL-3.0-only
package security

import (
	"net/http"
	"strings"
)

func BearerScope(v Validator, r *http.Request) (uint32, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return 0, ErrNoAuth
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, ErrBadAuth
	}
	return v.Validate(parts[1])
}
