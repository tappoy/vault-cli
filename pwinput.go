//go:build !mock
package main

import (
	"github.com/tappoy/pwinput"
)

// Normal implementation
func newPasswordInput() pwinput.PasswordInput {
	return pwinput.NewPasswordInput()
}

func setDummyPassword(dummy string) string {
	return ""
}

func setInterruptPassword() string {
	return ""
}
