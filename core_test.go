package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	testEnv = env{
		VaultDir:    "tmp/core/data",
		VaultLogDir: "tmp/core/log",
		VaultName:   "",
	}
)

func makeEnv(name string) env {
	return env{
		VaultDir:    testEnv.VaultDir,
		VaultLogDir: testEnv.VaultLogDir,
		VaultName:   name,
	}
}

func cleanCore(e env) {
	os.RemoveAll(filepath.Join(e.VaultDir, e.VaultName))
	os.RemoveAll(filepath.Join(e.VaultLogDir, e.VaultName))
}

func split(s string) []string {
	return strings.Split(s, " ")
}

func wName(name string) string {
	return "tmp/core/" + name + "_stdout.txt"
}

func makeStdout(t *testing.T, name string) *os.File {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		t.Fatal(err)
	}
	w, err := os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

func run(t *testing.T, o *option, want int) {
	t.Logf("o: %v", o.args)
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

func TestCore_Help(t *testing.T) {
	testName := "core_help"
	e := makeEnv(testName)
	cleanCore(e)
	w := makeStdout(t, wName(testName))
	defer w.Close()
	f := flags{name: nil}
	o := parse(e, f, split("vault-cli help"), w)
	run(t, o, 0)
}

func TestCore_Version(t *testing.T) {
	testName := "core_version"
	e := makeEnv(testName)
	cleanCore(e)
	w := makeStdout(t, wName(testName))
	defer w.Close()
	f := flags{name: nil}
	o := parse(e, f, split("vault-cli version"), w)
	run(t, o, 0)
}

func TestCore_Genpw(t *testing.T) {
	testName := "core_genpw"
	e := makeEnv(testName)
	cleanCore(e)
	w := makeStdout(t, wName(testName))
	defer w.Close()
	f := flags{name: nil}
	o := parse(e, f, split("vault-cli genpw"), w)
	run(t, o, 0)
}

func TestCore_Info(t *testing.T) {
	testName := "core_info"
	e := makeEnv(testName)
	cleanCore(e)
	w := makeStdout(t, wName(testName))
	defer w.Close()
	f := flags{name: nil}
	o := parse(e, f, split("vault-cli info"), w)
	run(t, o, 0)
}

func TestCore_Init(t *testing.T) {
	testName := "core_init"
	e := makeEnv(testName)
	cleanCore(e)

	{
		wn := wName(testName + "1")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli info"), w)
		run(t, o, 0)
		grepTrue(t, "init: false", wn)
	}

	{
		wn := wName(testName + "2")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli init"), w)
		run(t, o, 0)
	}

	{
		wn := wName(testName + "3")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli info"), w)
		run(t, o, 0)
		grepTrue(t, "init: true", wn)
	}
}

func setGetDelete(t *testing.T, testName string, dataDir bool, logDir bool) {
	e := makeEnv(testName)
	cleanCore(e)

	{
		wn := wName(testName + "1")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli init"), w)
		run(t, o, 0)
		grepTrue(t, "Init vault.", wn)

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
	}

	want := 0
	if !dataDir || !logDir {
		want = 1
		t.Logf("want: %v", want)
	}

	{
		wn := wName(testName + "2")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli set k1 k1value"), w)
		run(t, o, want)
	}

	{
		wn := wName(testName + "3")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli get k1"), w)
		run(t, o, want)
		if want == 0 {
			grepTrue(t, "k1value", wn)
		}
	}

	{
		wn := wName(testName + "4")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli delete k1"), w)
		run(t, o, want)
	}

	{
		wn := wName(testName + "5")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli get k1"), w)
		run(t, o, 1)
		if want == 0 {
			grepTrue(t, "Not found.", wn)
		}
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
	e := makeEnv(testName)
	cleanCore(e)

	{
		setDummyPassword("showtpw") // 7 characters
		wn := wName(testName + "1")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli init"), w)
		run(t, o, 1)
		grepTrue(t, "Wrong password.", wn)
	}

	{
		setInterruptPassword()
		wn := wName(testName + "2")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli init"), w)
		run(t, o, 1)
		grepTrue(t, "Interrupted.", wn)
	}

	{
		setDummyPassword("12345678") // valid password
		wn := wName(testName + "3")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli init"), w)
		run(t, o, 0)
		grepTrue(t, "Init vault.", wn)
	}

	{
		setDummyPassword("1234567890") // incorrect password
		wn := wName(testName + "4")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli set k1 k1value2"), w)
		run(t, o, 1)
		grepTrue(t, "Wrong password.", wn)
	}

}

func TestCore_UnknownCommand(t *testing.T) {
	testName := "core_unknown_command"
	e := makeEnv(testName)
	cleanCore(e)

	{
		wn := wName(testName + "1")
		w := makeStdout(t, wn)
		defer w.Close()
		f := flags{name: nil}
		o := parse(e, f, split("vault-cli unknown-command"), w)
		run(t, o, 1)
		grepTrue(t, "Unknown command. Run vault-cli help", wn)
	}
}
