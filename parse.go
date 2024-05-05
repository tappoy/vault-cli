package main

import (
	"flag" // TODO: To be replaced because the behavior is not good when extra flags are mixed in.
	"path/filepath"

	"github.com/tappoy/env"
)

func parse() *option {
	// get environment variables
	vaultDir := env.Getenv("VAULT_DIR", "/srv")
	vaultLogDir := env.Getenv("VAULT_LOG_DIR", "/var/log")
	vaultName := env.Getenv("VAULT_NAME", "vault")

	// make flag set
	flagset := flag.NewFlagSet("vault-cli", flag.ContinueOnError)
	flagset.SetOutput(env.Err)

	// parse flags
	var n string
	flagset.StringVar(&n, "n", "", "vault name")

	if flagset.Parse(env.Args[1:]) != nil {
		return nil
	}

	args := append(env.Args[:1], flagset.Args()...)
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
		w:        env.Out,
		args:     args,
	}
}
