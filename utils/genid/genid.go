// Package genid provides utilities for generating random ids.
package genid

import (
	"github.com/teris-io/shortid"
)

func Short() string {
	return shortid.MustGenerate()
}
