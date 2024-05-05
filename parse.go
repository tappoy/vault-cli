package main

import (
	"flag"
	"path/filepath"

	"github.com/tappoy/env"
)

func parse() *option {
	// get environment variables
	vaultDir := env.Getenv("VAULT_DIR", "/srv")
	vaultLogDir := env.Getenv("VAULT_LOG_DIR", "/var/log")
	vaultName := env.Getenv("VAULT_NAME", "vault")

	// make flag set
	flags := flag.NewFlagSet("vault-cli", flag.ContinueOnError)
	flags.SetOutput(env.Err)

	// parse flags
	var (
		n = flags.String("n", "", "vault name")
	)

	if flags.Parse(env.Args) != nil {
		return nil
	}

	args := flags.Args()

	var command string
	if len(args) < 2 {
		command = ""
	} else {
		command = args[1]
	}

	var name string
	if *n != "" {
		name = *n
	} else {
		name = vaultName
	}

	return &option{
		command:  command,
		name:     name,
		vaultDir: filepath.Join(vaultDir, name),
		logDir:   filepath.Join(vaultLogDir, name),
		w:        env.Out,
		args:     args,
	}
}
