//go:build test

package main

import (
	"github.com/tappoy/pwinput"
)

var dummyPassword = "dummyPassword"

func newPasswordInput() pwinput.PasswordInput {
	return pwinput.NewDummyPasswordInput(dummyPassword)
}
