// Package genid provides utilities for generating random ids.
package genid

import (
	"strings"

	"github.com/teris-io/shortid"
)

func Short() string {
	return strings.ToLower(shortid.MustGenerate())
}
