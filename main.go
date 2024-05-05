package main

import (
	"github.com/tappoy/env"
)

func main() {
	// parse arguments
	o := parse()
	if o == nil {
		env.Exit(1)
	}

	// run command
	env.Exit(o.run())
}
