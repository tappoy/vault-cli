//go:build test

package main

import (
	"fmt"
	"github.com/tappoy/pwinput"
)

var dummyPassword = "dummyPassword"

func newPasswordInput() pwinput.PasswordInput {
	fmt.Println("RUNNING MOCK: dummyPassword: ", dummyPassword)
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
