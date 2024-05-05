package main

import (
	"io"
	"testing"
)

type wants struct {
	command  string
	name     string
	vaultDir string
	logDir   string
	w        io.Writer
	args     []string
}

func check(t *testing.T, e env, f flags, args []string, w io.Writer, want wants) {
	opt := parse(e, f, args, w)
	if opt.command != want.command {
		t.Errorf("command: got %q, want %q", opt.command, want.command)
	}
	if opt.name != want.name {
		t.Errorf("name: got %q, want %q", opt.name, want.name)
	}
	if opt.vaultDir != want.vaultDir {
		t.Errorf("vaultDir: got %q, want %q", opt.vaultDir, want.vaultDir)
	}
	if opt.logDir != want.logDir {
		t.Errorf("logDir: got %q, want %q", opt.logDir, want.logDir)
	}
	if opt.w != want.w {
		t.Errorf("w: got %q, want %q", opt.w, want.w)
	}
}

func TestParse(t *testing.T) {
	e := env{
		VaultName:   "vault",
		VaultDir:    "/srv",
		VaultLogDir: "/var/log",
	}
	f := flags{
		name: nil,
	}
	args := []string{"arg1", "arg2", "arg3"}
	w := io.Discard
	want := wants{
		command:  "arg2",
		name:     "vault",
		vaultDir: "/srv/vault",
		logDir:   "/var/log/vault",
		w:        w,
		args:     args,
	}
	check(t, e, f, args, w, want)

	name := "joey"
	f = flags{name: &name}
	want = wants{
		command:  "arg2",
		name:     "joey",
		vaultDir: "/srv/joey",
		logDir:   "/var/log/joey",
		w:        w,
		args:     args,
	}
	check(t, e, f, args, w, want)

	args = []string{"arg1"}
	want = wants{
		command:  "",
		name:     "joey",
		vaultDir: "/srv/joey",
		logDir:   "/var/log/joey",
		w:        w,
		args:     args,
	}
	check(t, e, f, args, w, want)
}
