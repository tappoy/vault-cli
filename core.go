package main

import (
	"github.com/tappoy/env"
	"github.com/tappoy/logger"
	"github.com/tappoy/vault"
	ver "github.com/tappoy/version"

	"fmt"
	"os"
)

type option struct {
	command  string
	name     string
	vaultDir string
	logDir   string
	args     []string
}

// input password
func (o *option) inputPassword(logger *logger.Logger) (string, bool) {
	env.Errf("Password: ")
	password, err := env.InputPassword()
	env.Errf("\n")
	if err != nil {
		logger.Info(fmt.Sprintf("Cannot get password.\terror:%v", err))
		env.Errf("Interrupted.")
		return "", false
	}
	return password, true
}

// create vault
func (o *option) createVault(logger *logger.Logger) (*vault.Vault, bool) {
	password, ok := o.inputPassword(logger)
	if !ok {
		return nil, false
	}

	v, err := vault.NewVault(password, o.vaultDir)
	if err == nil {
		return v, true
	}

	switch err {
	case vault.ErrInvalidPasswordLength:
		msg := fmt.Sprintf("The password must be 8 to 32 characters.")
		env.Errf("%s\n", msg)
		logger.Notice(msg)
	case vault.ErrPasswordIncorrect:
		msg := fmt.Sprintf("Wrong password.")
		env.Errf("%s\n", msg)
		logger.Notice(msg)
	default: // TODO: cover. Permission.
		msg := fmt.Sprintf("Cannot open vault.\terror:%v\tvaultDir:%s", err, o.vaultDir)
		env.Errf("%s\n", msg)
		logger.Info(msg)
	}
	return nil, false
}

// get key
func (o *option) getKey() string {
	return o.args[2]
}

// create logger
func (o *option) createLogger() *logger.Logger {
	logger, err := logger.NewLogger(o.logDir)
	if err != nil { // TODO: cover. Permission.
		env.Errf("Cannot create logger.\terror:%v\tlogDir:%s\n", err, o.logDir)
		return nil
	}
	return logger
}

func (o *option) run() int {
	switch o.command {
	case "", "help":
		return o.usage()
	case "version":
		return o.version()
	case "genpw":
		return o.generatePassword()
	case "info":
		return o.info()
	case "init":
		return o.init()
	case "set":
		return o.set()
	case "read":
		return o.read()
	case "get":
		return o.get()
	case "delete":
		return o.delete()
	default:
		env.Errf("Unknown command.\n")
		runHelpMessage()
		return 1
	}
}

// check vault initialized
func (o *option) checkVaultInitialized(v *vault.Vault, logger *logger.Logger) int {
	if !v.IsInitialized() { // TODO: cover. Permission.
		msg := fmt.Sprintf("Vault is not initialized.\tvaultDir:%s", o.vaultDir)
		env.Errf("%s\n", msg)
		logger.Info(msg)
		return 1
	} else {
		return 0
	}
}

// print usage
func (o *option) usage() int {
	env.Outf(usage)
	return 0
}

// print version
func (o *option) version() int {
	env.Outf("vault-cli version %s\n", ver.Version())
	return 0
}

// print random password
func (o *option) generatePassword() int {
	password := vault.GeneratePassword()
	env.Outf(password)
	return 0
}

// print vault info
func (o *option) info() int {
	env.Outf("[Vault Info]\n")
	env.Outf("  name: %v\n", o.name)
	env.Outf("  data: %v\n", o.vaultDir)
	env.Outf("  log : %v\n", o.logDir)
	env.Outf("  init: %v\n", vault.IsInitialized(o.vaultDir))
	return 0
}

func (o *option) init() int {
	logger := o.createLogger()
	if logger == nil { // TODO: cover
		return 1
	}

	v, ok := o.createVault(logger)
	if !ok {
		return 1
	}

	err := v.Init()
	if err != nil { // TODO: cover
		msg := fmt.Sprintf("Cannot init vault.\terror:%v\tvaultDir:%s", err, o.vaultDir)
		env.Errf("%s\n", msg)
		logger.Notice(msg)
		return 1
	}

	msg := fmt.Sprintf("Init vault.\tvaultDir:%s", o.vaultDir)
	env.Outf(msg)
	logger.Notice(msg)
	return 0
}

func (o *option) set() int {
	logger := o.createLogger()
	if logger == nil {
		return 1
	}

	v, ok := o.createVault(logger)
	if !ok {
		return 1
	}

	if o.checkVaultInitialized(v, logger) != 0 {
		return 1
	}

	key := o.getKey()

	var value string
	if len(o.args) >= 4 {
		value = o.args[3]
	} else { // TODO: cover. No set value.
		value = ""
	}

	if err := v.Set(key, value); err != nil { // TODO: cover. Never happen?
		msg := fmt.Sprintf("Cannot set.\tkey:%s\terror:%v", key, err)
		env.Errf("%s\n", msg)
		logger.Info(msg)
		return 1
	}

	msg := fmt.Sprintf("set\tkey:%s", key)
	logger.Info(msg)
	env.Outf("Set successfully.\n")
	return 0
}

func (o *option) read() int {
	logger := o.createLogger()
	if logger == nil {
		return 1
	}

	v, ok := o.createVault(logger)
	if !ok {
		return 1
	}

	if o.checkVaultInitialized(v, logger) != 0 {
		return 1
	}

	key := o.getKey()

	var file string
	if len(o.args) >= 4 {
		file = o.args[3]
	} else { // TODO: cover. No set file.
		file = ""
	}

	bytes, err := os.ReadFile(file)
	if err != nil {
		msg := fmt.Sprintf("Cannot read file.\tfile:%s\terror:%v", file, err)
		env.Errf("%s\n", msg)
		logger.Info(msg)
		return 1
	}

	value := string(bytes)
	if err := v.Set(key, value); err != nil { // TODO: cover. Never happen?
		msg := fmt.Sprintf("Cannot set.\tkey:%s\terror:%v", key, err)
		env.Errf("%s\n", msg)
		logger.Info(msg)
		return 1
	}

	msg := fmt.Sprintf("read\tkey:%s", key)
	logger.Info(msg)
	env.Outf("Read successfully.\n")
	return 0
}

func (o *option) get() int {
	logger := o.createLogger()
	if logger == nil {
		return 1
	}

	v, ok := o.createVault(logger)
	if !ok { // TODO: cover
		return 1
	}

	if o.checkVaultInitialized(v, logger) != 0 {
		return 1
	}

	key := o.getKey()

	value, err := v.Get(key)
	if err != nil {
		switch err {
		case vault.ErrKeyNotFound:
			msg := fmt.Sprintf("Not found.\tkey:%s\n", key)
			env.Errf("%s\n", msg)
			logger.Info(msg)
		default: // TODO: cover. Never happen?
			msg := fmt.Sprintf("Cannot get.\tkey:%s error:%v", key, err)
			env.Errf("%s\n", msg)
			logger.Info(msg)
		}
		return 1
	}

	msg := fmt.Sprintf("get\tkey:%s", key)
	logger.Info(msg)
	env.Outf("%s\n", value)
	return 0
}

func (o *option) delete() int {
	logger := o.createLogger()
	if logger == nil {
		return 1
	}

	v, ok := o.createVault(logger)
	if !ok { // TODO: cover
		return 1
	}

	if o.checkVaultInitialized(v, logger) != 0 {
		return 1
	}

	key := o.getKey()

	if err := v.Delete(key); err != nil { // TODO: cover
		switch err {
		case vault.ErrKeyNotFound: // TODO: cover
			msg := fmt.Sprintf("Not found.\tkey:%s", key)
			env.Errf("%s\n", msg)
			logger.Info(msg)
		default: // TODO: cover. Never happen?
			msg := fmt.Sprintf("Cannot delete.\tkey:%s error:%v", key, err)
			env.Errf("%s\n", msg)
			logger.Info(msg)
		}
		return 1 // TODO: cover
	}

	msg := fmt.Sprintf("delete\tkey:%s", key)
	logger.Info(msg)
	env.Outf("%s is deleted.\n", key)
	return 0
}
