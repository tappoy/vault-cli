package main

import (
	"flag"
	te "github.com/tappoy/env"
	"os"
)

type env struct {
	VaultDir    string
	VaultLogDir string
	VaultName   string
}

type flags struct {
	name *string
}

func main() {
	// get environment variables
	e := env{
		VaultDir:    te.GetEnvString("VAULT_DIR", "/srv"),
		VaultLogDir: te.GetEnvString("VAULT_LOG_DIR", "/var/log"),
		VaultName:   te.GetEnvString("VAULT_NAME", "vault"),
	}

	// parse flags
	var (
		n = flag.String("n", "vault", "vault name")
	)
	flag.Parse()
	// flag.Args() does not include $0 so we need to append it.
	args := append(os.Args[:1], flag.Args()...)

	// parse arguments
	o := parse(e, flags{name: n}, args, os.Stdout)

	// run command
	os.Exit(o.run())
}
