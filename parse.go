package main

import (
	"github.com/tappoy/env"
	"github.com/tappoy/flag"

	"path/filepath"
)

func parse() *option {
	// get environment variables
	vaultDir := env.Getenv("VAULT_DIR", "/srv")
	vaultLogDir := env.Getenv("VAULT_LOG_DIR", "/var/log")
	vaultName := env.Getenv("VAULT_NAME", "vault")

	// make flag set
	flagset := flag.NewFlagSet(env.Args)

	// parse flags
	var n string
	flagset.StringVar(&n, "n", "")

	err := flagset.Parse()
	if err != nil {
		parseErr := err.(*flag.ParseError)
		switch parseErr.Err {
		case flag.ErrInvalidValue:
			env.Errf("Invalid value for flag %s: %s\n", parseErr.Arg, parseErr.Value)
		}
		runHelpMessage()
		return nil
	}

	args := flagset.Args()
	var command string
	if len(args) < 2 {
		command = ""
	} else {
		command = args[1]
	}

	var name string
	if n != "" {
		name = n
	} else {
		name = vaultName
	}

	return &option{
		command:  command,
		name:     name,
		vaultDir: filepath.Join(vaultDir, name),
		logDir:   filepath.Join(vaultLogDir, name),
		args:     args,
	}
}
