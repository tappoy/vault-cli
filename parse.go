package main

import (
	"fmt"
	"github.com/tappoy/logger"
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

func newOptions(e env, args []string, w io.Writer) (*option, int) {
	if len(args) < 2 {
		args = append(args, "help")
	}

	name := ""

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
	}

	logDir := filepath.Join(getLogDirRoot(e), name)
	logger, err := logger.NewLogger(logDir)
	if err != nil {
		fmt.Fprintf(w, "Cannot create logger.\tlogDir:%s\terror:%v\n", logDir, err)
		return nil, 1
	}

	return &option{
		command:      command,
		name:         name,
		password:     "",
		vaultDirRoot: getVaultDirRoot(e),
		logger:       logger,
		w:            w,
		args:         args,
	}, 0
}
