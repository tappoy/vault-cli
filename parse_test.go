package main

import (
	"github.com/tappoy/env"
	"reflect"
	"testing"
)

func want(command, name, vaultDir, logDir string, args []string) *option {
	return &option{command: command, name: name, vaultDir: vaultDir, logDir: logDir, w: env.Out, args: args}
}

// TestParse tests the parse function.
func TestParse(t *testing.T) {
	// define nil environment
	env0 := env.Env{"VAULT_DIR": "", "VAULT_LOG_DIR": "", "VAULT_NAME": ""}
	env1 := env.Env{"VAULT_DIR": "tmp/parse/data", "VAULT_LOG_DIR": "tmp/parse/log", "VAULT_NAME": "parse_test"}

	// test cases
	cases := []struct {
		args []string
		envs env.Env
		want *option
	}{
		{
			args: split("vault-cli help"),
			envs: env0,
			want: want("help", "vault", "/srv/vault", "/var/log/vault", split("vault-cli help")),
		},
		{
			args: split("vault-cli help"),
			envs: env1,
			want: want("help", "parse_test", "tmp/parse/data/parse_test", "tmp/parse/log/parse_test", split("vault-cli help")),
		},
	}

	// run tests
	for _, c := range cases {
		env.DummyEnv = c.envs
		env.Args = c.args
		got := parse()
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ERROR:\nargs %s\ngot  %v\nwant %v\n", c.args, got, c.want)
		}
	}
}
