package main

import (
	_ "embed"
	"github.com/tappoy/env"
)

//go:embed Usage.txt
var usage string

var (
	VaultDir    = "/srv"
	VaultLogDir = "/var/log"
	VaultName   = "vault"
)

func runHelpMessage() {
	env.Errf("Run %s help\n", env.Args[0])
}

func main() {
	// parse arguments
	o := parse()
	if o == nil {
		env.Exit(1)
	}

	// run command
	env.Exit(o.run())
}
