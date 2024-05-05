package main

import (
	"io"
	"path/filepath"
)

func parse(e env, f flags, args []string, w io.Writer) *option {
	var command string
	if len(args) < 2 {
		command = ""
	} else {
		command = args[1]
	}

	var name string
	if f.name != nil {
		name = *f.name
	} else {
		name = e.VaultName
	}

	return &option{
		command:  command,
		name:     name,
		vaultDir: filepath.Join(e.VaultDir, name),
		logDir:   filepath.Join(e.VaultLogDir, name),
		w:        w,
		args:     args,
	}
}
