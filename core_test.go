package main

import (
	"github.com/tappoy/env"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setEnv(name string) {
	env.DummyEnv = env.Env{
		"VAULT_DIR":     testRoot + "/core/data",
		"VAULT_LOG_DIR": testRoot + "/core/log",
		"VAULT_NAME":    name,
	}
}

func wName(name string) string {
	return testRoot + "/core/" + name + "_stdout.txt"
}

func wNameErr(name string) string {
	return testRoot + "/core/" + name + "_stderr.txt"
}

func setStdout(t *testing.T, stdout, stderr string) {
	if err := os.MkdirAll(filepath.Dir(stdout), 0755); err != nil {
		t.Fatal(err)
	}
	o, err := os.Create(stdout)
	if err != nil {
		t.Fatal(err)
	}
	env.Out = o
	e, err := os.Create(stderr)
	if err != nil {
		t.Fatal(err)
	}
	env.Err = e
}

func run(t *testing.T, o *option, want int) {
	if rc := o.run(); rc != want {
		t.Errorf("ERROR: got %v, want %v", rc, want)
		t.Logf("o: %v", o)
	}
}

func grepFalse(t *testing.T, search string, file string) {
	grep(t, search, file, false)
}

func grepTrue(t *testing.T, search string, file string) {
	grep(t, search, file, true)
}

func grep(t *testing.T, search string, file string, want bool) {
	f, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// read file as a string
	output, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(output), search) != want {
		t.Errorf("Unmatched %v: [%v]", want, search)
	}
}

func doTest(t *testing.T, testName, wsuffix, args string, want int, search, searchErr string) (*option, string, string) {
	wn := wName(testName + wsuffix)
	we := wNameErr(testName + wsuffix)
	setStdout(t, wn, we)
	setEnv(testName)
	env.Args = split(args)
	o := parse()
	run(t, o, want)
	if search != "" {
		grepTrue(t, search, wn)
	}
	if searchErr != "" {
		grepTrue(t, searchErr, we)
	}
	return o, wn, we
}

func TestCore_Help(t *testing.T) {
	doTest(t, "core_help", "", "vault-cli help", 0, "Usage:", "")
}

func TestCore_Version(t *testing.T) {
	doTest(t, "core_version", "", "vault-cli version", 0, "version", "")
}

func TestCore_Genpw(t *testing.T) {
	doTest(t, "core_genpw", "", "vault-cli genpw", 0, "", "")
}

func TestCore_Info(t *testing.T) {
	doTest(t, "core_info", "", "vault-cli info", 0, "init: false", "")
}

func TestCore_Init(t *testing.T) {
	testName := "core_init"
	doTest(t, testName, "1", "vault-cli info", 0, "init: false", "")
	doTest(t, testName, "2", "vault-cli init", 0, "Init vault.", "")
	doTest(t, testName, "3", "vault-cli info", 0, "init: true", "")
}

func setGetDelete(t *testing.T, testName string, dataDir bool, logDir bool) {

	o, wn, we := doTest(t, testName, "1", "vault-cli init", 0, "Init vault.", "")

	if !dataDir {
		// set read only to data dir
		os.Chmod(o.vaultDir, 0400)
		defer os.Chmod(o.vaultDir, 0700)
	}
	if !logDir {
		// set read only to log dir
		os.Chmod(o.logDir, 0400)
		defer os.Chmod(o.logDir, 0700)
	}

	want := 0
	if !dataDir || !logDir {
		want = 1
		t.Logf("want: %v", want)
	}

	o, wn, we = doTest(t, testName, "2", "vault-cli set k1 k1value", want, "", "")
	o, wn, we = doTest(t, testName, "3", "vault-cli get k1", want, "", "")
	if want == 0 {
		grepTrue(t, "k1value", wn)
	} else {
		grepTrue(t, ".", we)
	}

	o, wn, we = doTest(t, testName, "4", "vault-cli delete k1", want, "", "")
	o, wn, we = doTest(t, testName, "5", "vault-cli get k1", 1, "", "")
	if want == 0 {
		grepTrue(t, "Not found.", we)
	} else {
		grepTrue(t, ".", we)
	}
}

func TestCore_SetGetDelete(t *testing.T) {
	testName := "core_set_get_delete"
	setGetDelete(t, testName, true, true)
}

func TestCore_SetGetDeleteWithReadOnlyDataDir(t *testing.T) {
	testName := "core_set_get_delete_with_read_only_data_dir"
	setGetDelete(t, testName, false, true)
}

func TestCore_SetGetDeleteWithReadOnlyLogDir(t *testing.T) {
	testName := "core_set_get_delete_with_read_only_log_dir"
	setGetDelete(t, testName, true, false)
}

func TestCore_PasswordIncorrect(t *testing.T) {
	testName := "core_password_incorrect"

	setDummyPassword("showtpw") // 7 characters
	doTest(t, testName, "1", "vault-cli init", 1, "", "The password must be 8 to 32 characters.")

	setInterruptPassword()
	doTest(t, testName, "2", "vault-cli init", 1, "", "Interrupted.")

	setDummyPassword("12345678") // valid password
	doTest(t, testName, "3", "vault-cli init", 0, "Init vault.", "")

	setDummyPassword("1234567890") // incorrect password
	doTest(t, testName, "4", "vault-cli set k1 k1value2", 1, "", "Wrong password.")
}

func TestCore_UnknownCommand(t *testing.T) {
	testName := "core_unknown_command"
	doTest(t, testName, "1", "vault-cli unknown", 1, "", "Unknown command.")
}
