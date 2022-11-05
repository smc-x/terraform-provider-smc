// Package logging provides logging utilities.
package logging

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var (
	ErrDiags = errors.New("diags has error")
)

func WarnIf(verb string, err error, body diag.Diagnostics) {
	if err != nil {
		body.Append(diag.NewWarningDiagnostic(
			"failed to "+verb,
			err.Error(),
		))
	}
}

func PanicIf(verb string, err error, body diag.Diagnostics) {
	if err != nil {
		body.Append(diag.NewErrorDiagnostic(
			"failed to "+verb,
			err.Error(),
		))
		panic(ErrDiags)
	}
}

func PanicIfDiags(err, body diag.Diagnostics) {
	body.Append(err...)
	if body.HasError() {
		panic(ErrDiags)
	}
}

func RecoverDiags() {
	if e := recover(); e != nil && e != ErrDiags {
		panic(e)
	}
}
