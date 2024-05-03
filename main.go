package main

import (
	"os"
)

func main() {
	// parse environment variables
	e := env{
		VaultDir:    os.Getenv("VAULT_DIR"),
		VaultLogDir: os.Getenv("VAULT_LOG_DIR"),
		VaultName:   os.Getenv("VAULT_NAME"),
	}

	// parse arguments
	o, rc := newOptions(e, os.Args, os.Stdout)
	if o == nil || rc != 0 {
		os.Exit(rc)
	}

	os.Exit(o.run())
}
