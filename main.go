package main

import (
	_ "embed"
	"github.com/tappoy/env"
)

//go:embed Usage.txt
var usage string

func main() {
	// parse arguments
	o := parse()
	if o == nil {
		env.Exit(1)
	}

	// run command
	env.Exit(o.run())
}
