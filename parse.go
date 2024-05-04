package main

import (
	"io"
	"path/filepath"
	"strings"
)

func getName(nameIndex int, args []string, e env) string {
	if len(args) > nameIndex {
		return args[nameIndex]
	} else if name := e.VaultName; strings.TrimSpace(name) != "" {
		return strings.TrimSpace(name)
	} else {
		return "vault"
	}
}

func getVaultDirRoot(e env) string {
	if dir := e.VaultDir; dir != "" {
		return dir
	} else {
		return "/srv"
	}
}

func getLogDirRoot(e env) string {
	if dir := e.VaultLogDir; dir != "" {
		return dir
	} else {
		return "/var/log"
	}
}

type env struct {
	VaultDir    string
	VaultLogDir string
	VaultName   string
}

func parse(e env, args []string, w io.Writer) (*option, int) {
	if len(args) < 2 {
		args = append(args, "help")
	}

	var name string

	command := args[1]
	switch command {
	case "init", "info":
		name = getName(2, args, e)
	case "set":
		name = getName(4, args, e)
	case "get":
		name = getName(3, args, e)
	case "delete":
		name = getName(3, args, e)
	default:
		name = "vault"
	}

	return &option{
		command:  command,
		name:     name,
		vaultDir: filepath.Join(getVaultDirRoot(e), name),
		logDir:   filepath.Join(getLogDirRoot(e), name),
		w:        w,
		args:     args,
	}, 0
}
