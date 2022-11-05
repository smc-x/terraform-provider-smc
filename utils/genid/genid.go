// Package genid provides utilities for generating random ids.
package genid

import (
	"time"

	"github.com/teris-io/shortid"
)

var (
	sid = shortid.MustNew(0, "0123456789abcdefghijklmnopqrstuvwxyz", uint64(time.Now().Unix()))
)

func Short() string {
	return sid.MustGenerate()
}
