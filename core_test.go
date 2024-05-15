package main

import (
	"github.com/tappoy/env"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setGlobal(name string) {
	VaultDir = testRoot + "/core/data"
	VaultLogDir = testRoot + "/core/log"
	VaultName = name
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
	setGlobal(testName)
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
	env.SetDummyPassword("dummyPassword")
	doTest(t, "core_help", "", "vault-cli help", 0, "Usage:", "")
}

func TestCore_Version(t *testing.T) {
	env.SetDummyPassword("dummyPassword")
	doTest(t, "core_version", "", "vault-cli version", 0, "version", "")
}

func TestCore_Genpw(t *testing.T) {
	env.SetDummyPassword("dummyPassword")
	doTest(t, "core_genpw", "", "vault-cli genpw", 0, "", "")
}

func TestCore_Info(t *testing.T) {
	env.SetDummyPassword("dummyPassword")
	doTest(t, "core_info", "", "vault-cli info", 0, "init: false", "")
}

func TestCore_Init(t *testing.T) {
	env.SetDummyPassword("dummyPassword")
	testName := "core_init"
	doTest(t, testName, "1", "vault-cli info", 0, "init: false", "")
	doTest(t, testName, "2", "vault-cli init", 0, "Init vault.", "")
	doTest(t, testName, "3", "vault-cli info", 0, "init: true", "")
}

func TestCore_Read(t *testing.T) {
	env.SetDummyPassword("dummyPassword")
	testName := "core_read"
	os.WriteFile(testRoot+"/backup1", []byte("backup1\nuser\tjon\npass\tdoe\n"), 0644)
	doTest(t, testName, "1", "vault-cli init", 0, "Init vault.", "")
	doTest(t, testName, "2", "vault-cli read backup1 "+testRoot+"/backup1", 0, "Read successfully.", "")
	doTest(t, testName, "3", "vault-cli get backup1", 0, "pass\tdoe", "")
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
	env.SetDummyPassword("dummyPassword")
	testName := "core_set_get_delete"
	setGetDelete(t, testName, true, true)
}

func TestCore_SetGetDeleteWithReadOnlyDataDir(t *testing.T) {
	env.SetDummyPassword("dummyPassword")
	testName := "core_set_get_delete_with_read_only_data_dir"
	setGetDelete(t, testName, false, true)
}

func TestCore_SetGetDeleteWithReadOnlyLogDir(t *testing.T) {
	env.SetDummyPassword("dummyPassword")
	testName := "core_set_get_delete_with_read_only_log_dir"
	setGetDelete(t, testName, true, false)
}

func TestCore_PasswordIncorrect(t *testing.T) {
	testName := "core_password_incorrect"

	env.SetDummyPassword("showtpw") // 7 characters
	doTest(t, testName, "1", "vault-cli init", 1, "", "The password must be 8 to 32 characters.")

	env.SetDummyPassword(env.Interrupt)
	doTest(t, testName, "2", "vault-cli init", 1, "", "Interrupted.")

	env.SetDummyPassword("12345678") // valid password
	doTest(t, testName, "3", "vault-cli init", 0, "Init vault.", "")

	env.SetDummyPassword("1234567890") // incorrect password
	doTest(t, testName, "4", "vault-cli set k1 k1value2", 1, "", "Wrong password.")
}

func TestCore_UnknownCommand(t *testing.T) {
	env.SetDummyPassword("dummyPassword")
	testName := "core_unknown_command"
	doTest(t, testName, "1", "vault-cli unknown", 1, "", "Unknown command.")
}
