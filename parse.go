package main

import (
	"github.com/tappoy/env"
	"path/filepath"
)

func parse() *option {
	args := env.Args
	var command string
	if len(args) < 2 {
		command = ""
	} else {
		command = args[1]
	}

	return &option{
		command:  command,
		name:     VaultName,
		vaultDir: filepath.Join(VaultDir, VaultName),
		logDir:   filepath.Join(VaultLogDir, VaultName),
		args:     args,
	}
}
