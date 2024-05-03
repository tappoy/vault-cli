package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

var (
	testEnv = env{
		VaultDir:    "tmp/test/data",
		VaultLogDir: "tmp/test/log",
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
	os.RemoveAll(e.VaultDir)
	os.RemoveAll(e.VaultLogDir)
}

func split(s string) []string {
	return strings.Split(s, " ")
}

func wName(name string) string {
	return "tmp/" + name + "_stdout.txt"
}

func makeStdout(t *testing.T, name string) *os.File {
	w, err := os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

func run(t *testing.T, o *option, want int) {
	if rc := o.run(); rc != want {
		t.Errorf("ERROR: got %v, want %v", rc, want)
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
	o, _ := parse(e, split("vault-cli help"), w)
	run(t, o, 0)
}

func TestCore_Version(t *testing.T) {
	testName := "core_version"
	e := makeEnv(testName)
	cleanCore(e)
	w := makeStdout(t, wName(testName))
	defer w.Close()
	o, _ := parse(e, split("vault-cli version"), w)
	run(t, o, 0)
}

func TestCore_Genpw(t *testing.T) {
	testName := "core_genpw"
	e := makeEnv(testName)
	cleanCore(e)
	w := makeStdout(t, wName(testName))
	defer w.Close()
	o, _ := parse(e, split("vault-cli genpw"), w)
	run(t, o, 0)
}

func TestCore_Info(t *testing.T) {
	testName := "core_info"
	e := makeEnv(testName)
	cleanCore(e)
	w := makeStdout(t, wName(testName))
	defer w.Close()
	o, _ := parse(e, split("vault-cli info"), w)
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
		o, _ := parse(e, split("vault-cli info"), w)
		run(t, o, 0)
		grepTrue(t, "init: false", wn)
	}

	{
		wn := wName(testName + "2")
		w := makeStdout(t, wn)
		defer w.Close()
		o, _ := parse(e, split("vault-cli init"), w)
		run(t, o, 0)
	}

	{
		wn := wName(testName + "3")
		w := makeStdout(t, wn)
		defer w.Close()
		o, _ := parse(e, split("vault-cli info"), w)
		run(t, o, 0)
		grepTrue(t, "init: true", wn)
	}
}

func TestCore_SetGetDelete(t *testing.T) {
	testName := "core_set_get_delete"
	e := makeEnv(testName)
	cleanCore(e)

	{
		wn := wName(testName + "1")
		w := makeStdout(t, wn)
		defer w.Close()
		o, _ := parse(e, split("vault-cli init"), w)
		run(t, o, 0)
		grepTrue(t, "Init vault.", wn)
	}

	{
		wn := wName(testName + "2")
		w := makeStdout(t, wn)
		defer w.Close()
		o, _ := parse(e, split("vault-cli set k1 k1value"), w)
		run(t, o, 0)
	}

	{
		wn := wName(testName + "3")
		w := makeStdout(t, wn)
		defer w.Close()
		o, _ := parse(e, split("vault-cli get k1"), w)
		run(t, o, 0)
		grepTrue(t, "k1value", wn)
	}

	{
		wn := wName(testName + "4")
		w := makeStdout(t, wn)
		defer w.Close()
		o, _ := parse(e, split("vault-cli delete k1"), w)
		run(t, o, 0)
	}

	{
		wn := wName(testName + "5")
		w := makeStdout(t, wn)
		defer w.Close()
		o, _ := parse(e, split("vault-cli get k1"), w)
		run(t, o, 1)
		grepTrue(t, "Not found.", wn)
	}

}
