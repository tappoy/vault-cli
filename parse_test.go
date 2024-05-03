package main

import (
	"bytes"
	"strings"
	"testing"
)

type parseTest struct {
	command string
	e       env
	args    []string
	want    *option
	wantInt int
}

func makeParseTest(command string, e env, wantInt int, want *option) parseTest {
	return parseTest{
		command: command,
		e:       e,
		args:    strings.Split(command, " "),
		want:    want,
		wantInt: wantInt,
	}
}

func checkParse(t *testing.T, tests []parseTest) {
	w := new(bytes.Buffer)
	for _, test := range tests {
		t.Logf("Test: %s", test.command)
		got, gotInt := parse(test.e, test.args, w)
		if gotInt != test.wantInt {
			t.Errorf("ERROR: got %v, want %v", gotInt, test.wantInt)
		}
		if got.command != test.want.command {
			t.Errorf("ERROR: got %v, want %v", got.command, test.want.command)
		}
		if got.name != test.want.name {
			t.Errorf("ERROR: got %v, want %v", got.name, test.want.name)
		}
		if got.logDir != test.want.logDir {
			t.Errorf("ERROR: got %v, want %v", got.logDir, test.want.logDir)
		}
		if got.vaultDir != test.want.vaultDir {
			t.Errorf("ERROR: got %v, want %v", got.vaultDir, test.want.vaultDir)
		}
	}
}

func TestNewOptionsWithNilEnvNilName(t *testing.T) {
	testEnv := env{
		VaultDir:    "",
		VaultLogDir: "",
		VaultName:   "",
	}

	tests := []parseTest{
		makeParseTest("vault-cli help", testEnv, 0, &option{
			command:  "help",
			name:     "vault",
			logDir:   "/var/log/vault",
			vaultDir: "/srv/vault",
		}),
		makeParseTest("vault-cli info", testEnv, 0, &option{
			command:  "info",
			name:     "vault",
			logDir:   "/var/log/vault",
			vaultDir: "/srv/vault",
		}),
		makeParseTest("vault-cli init", testEnv, 0, &option{
			command:  "init",
			name:     "vault",
			logDir:   "/var/log/vault",
			vaultDir: "/srv/vault",
		}),
		makeParseTest("vault-cli set key val", testEnv, 0, &option{
			command:  "set",
			name:     "vault",
			logDir:   "/var/log/vault",
			vaultDir: "/srv/vault",
		}),
		makeParseTest("vault-cli get key", testEnv, 0, &option{
			command:  "get",
			name:     "vault",
			logDir:   "/var/log/vault",
			vaultDir: "/srv/vault",
		}),
		makeParseTest("vault-cli delete key", testEnv, 0, &option{
			command:  "delete",
			name:     "vault",
			logDir:   "/var/log/vault",
			vaultDir: "/srv/vault",
		}),
	}

	checkParse(t, tests)
}

func TestNewOptionsWithNilEnvWithName(t *testing.T) {
	name := "test"

	testEnv := env{
		VaultDir:    "",
		VaultLogDir: "",
		VaultName:   "",
	}

	tests := []parseTest{
		makeParseTest("vault-cli help "+name, testEnv, 0, &option{
			// TODO: name should be got from the -name flag
			command:  "help",
			name:     "vault",
			logDir:   "/var/log/" + "vault",
			vaultDir: "/srv/" + "vault",
		}),
		makeParseTest("vault-cli info "+name, testEnv, 0, &option{
			command:  "info",
			name:     name,
			logDir:   "/var/log/" + name,
			vaultDir: "/srv/" + name,
		}),
		makeParseTest("vault-cli init "+name, testEnv, 0, &option{
			command:  "init",
			name:     name,
			logDir:   "/var/log/" + name,
			vaultDir: "/srv/" + name,
		}),
		makeParseTest("vault-cli set key val "+name, testEnv, 0, &option{
			command:  "set",
			name:     name,
			logDir:   "/var/log/" + name,
			vaultDir: "/srv/" + name,
		}),
		makeParseTest("vault-cli get key "+name, testEnv, 0, &option{
			command:  "get",
			name:     name,
			logDir:   "/var/log/" + name,
			vaultDir: "/srv/" + name,
		}),
		makeParseTest("vault-cli delete key "+name, testEnv, 0, &option{
			command:  "delete",
			name:     name,
			logDir:   "/var/log/" + name,
			vaultDir: "/srv/" + name,
		}),
	}

	checkParse(t, tests)
}

func TestNewOptionsWithEnvNilName(t *testing.T) {
	testEnv := env{
		VaultDir:    "/env",
		VaultLogDir: "/env/log",
		VaultName:   "vne",
	}

	tests := []parseTest{
		makeParseTest("vault-cli help", testEnv, 0, &option{
			command:  "help",
			name:     "vault", // TODO: name should be got from the -name flag
			logDir:   "/env/log/vault",
			vaultDir: "/env/vault",
		}),
		makeParseTest("vault-cli info", testEnv, 0, &option{
			command:  "info",
			name:     "vne",
			logDir:   "/env/log/vne",
			vaultDir: "/env/vne",
		}),
		makeParseTest("vault-cli init", testEnv, 0, &option{
			command:  "init",
			name:     "vne",
			logDir:   "/env/log/vne",
			vaultDir: "/env/vne",
		}),
		makeParseTest("vault-cli set key val", testEnv, 0, &option{
			command:  "set",
			name:     "vne",
			logDir:   "/env/log/vne",
			vaultDir: "/env/vne",
		}),
		makeParseTest("vault-cli get key", testEnv, 0, &option{
			command:  "get",
			name:     "vne",
			logDir:   "/env/log/vne",
			vaultDir: "/env/vne",
		}),
		makeParseTest("vault-cli delete key", testEnv, 0, &option{
			command:  "delete",
			name:     "vne",
			logDir:   "/env/log/vne",
			vaultDir: "/env/vne",
		}),
	}

	checkParse(t, tests)
}

func TestNewOptionsWithEnvWithName(t *testing.T) {
	name := "vne2"

	testEnv := env{
		VaultDir:    "/env",
		VaultLogDir: "/env/log",
		VaultName:   "this_is_not_used",
	}

	tests := []parseTest{
		makeParseTest("vault-cli help "+name, testEnv, 0, &option{
			command:  "help",
			name:     "vault", // TODO: name should be got from the -name flag
			logDir:   "/env/log/" + "vault",
			vaultDir: "/env/" + "vault",
		}),
		makeParseTest("vault-cli info "+name, testEnv, 0, &option{
			command:  "info",
			name:     name,
			logDir:   "/env/log/" + name,
			vaultDir: "/env/" + name,
		}),
		makeParseTest("vault-cli init "+name, testEnv, 0, &option{
			command:  "init",
			name:     name,
			logDir:   "/env/log/" + name,
			vaultDir: "/env/" + name,
		}),
		makeParseTest("vault-cli set key val "+name, testEnv, 0, &option{
			command:  "set",
			name:     name,
			logDir:   "/env/log/" + name,
			vaultDir: "/env/" + name,
		}),
		makeParseTest("vault-cli get key "+name, testEnv, 0, &option{
			command:  "get",
			name:     name,
			logDir:   "/env/log/" + name,
			vaultDir: "/env/" + name,
		}),
		makeParseTest("vault-cli delete key "+name, testEnv, 0, &option{
			command:  "delete",
			name:     name,
			logDir:   "/env/log/" + name,
			vaultDir: "/env/" + name,
		}),
	}

	checkParse(t, tests)
}
