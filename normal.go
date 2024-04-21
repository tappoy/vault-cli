//go:build !test

package main

import (
	"github.com/tappoy/pwinput"
)

// Normal implementation
func newPasswordInput() pwinput.PasswordInput {
	return pwinput.NewPasswordInput()
}
