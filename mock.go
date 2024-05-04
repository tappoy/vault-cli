//go:build test

package main

import (
	"github.com/tappoy/pwinput"
)

var dummyPassword = "dummyPassword"

func newPasswordInput() pwinput.PasswordInput {
	return pwinput.NewDummyPasswordInput(dummyPassword)
}

func setDummyPassword(dummy string) string {
	dummyPassword = dummy
	return dummyPassword
}

func setInterruptPassword() string {
	dummyPassword = pwinput.Interrupt
	return dummyPassword
}
