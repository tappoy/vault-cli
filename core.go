package main

import (
	"github.com/tappoy/logger"
	"github.com/tappoy/vault"
	ver "github.com/tappoy/version"

	"fmt"
	"io"
)

type option struct {
	command  string
	name     string
	vaultDir string
	logDir   string
	w        io.Writer
	args     []string
}

// input password
func (o *option) inputPassword(logger *logger.Logger) (string, bool) {
	fmt.Fprint(o.w, "Password: ")
	pwi := newPasswordInput()
	password, err := pwi.InputPassword()
	fmt.Fprintln(o.w)
	if err != nil {
		msg := fmt.Sprintf("Cannot get password.\terror:%v", err)
		logger.Info(msg)
		fmt.Fprintf(o.w, "%v.", err)
		fmt.Fprintln(o.w, "Interrupted.")
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
	case vault.ErrInvalidPasswordLength, vault.ErrPasswordIncorrect:
		msg := fmt.Sprintf("Wrong password.")
		fmt.Fprintln(o.w, msg)
		logger.Notice(msg)
	default: // TODO: cover. Permission.
		msg := fmt.Sprintf("Cannot open vault.\terror:%v\tvaultDir:%s", err, o.vaultDir)
		fmt.Fprintln(o.w, msg)
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
		fmt.Fprintf(o.w, "Cannot create logger.\terror:%v\tlogDir:%s\n", err, o.logDir)
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
	case "get":
		return o.get()
	case "delete":
		return o.delete()
	default:
		fmt.Fprintf(o.w, "Unknown command.")
		runHelpMessage()
		return 1
	}
}

// check vault initialized
func (o *option) checkVaultInitialized(v *vault.Vault, logger *logger.Logger) int {
	if !v.IsInitialized() { // TODO: cover. Permission.
		msg := fmt.Sprintf("Vault is not initialized.\tvaultDir:%s", o.vaultDir)
		fmt.Fprintln(o.w, msg)
		logger.Info(msg)
		return 1
	} else {
		return 0
	}
}

// print usage
func (o *option) usage() int {
	fmt.Fprintf(o.w, usage)
	return 0
}

// print version
func (o *option) version() int {
	fmt.Fprintf(o.w, "vault-cli version %s\n", ver.Version())
	return 0
}

// print random password
func (o *option) generatePassword() int {
	password := vault.GeneratePassword()
	fmt.Fprintln(o.w, password)
	return 0
}

// print vault info
func (o *option) info() int {
	fmt.Fprintln(o.w, "[Vault Info]")
	fmt.Fprintln(o.w, "  name:", o.name)
	fmt.Fprintln(o.w, "  data:", o.vaultDir)
	fmt.Fprintln(o.w, "  log :", o.logDir)
	fmt.Fprintln(o.w, "  init:", vault.IsInitialized(o.vaultDir))
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
		fmt.Fprintln(o.w, msg)
		logger.Notice(msg)
		return 1
	}

	msg := fmt.Sprintf("Init vault.\tvaultDir:%s", o.vaultDir)
	fmt.Fprintln(o.w, msg)
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
		fmt.Fprintln(o.w, msg)
		logger.Info(msg)
		return 1
	}

	msg := fmt.Sprintf("set\tkey:%s", key)
	logger.Info(msg)
	fmt.Fprintln(o.w, "Set successfully.")
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
			msg := fmt.Sprintf("Not found.\tkey:%s", key)
			fmt.Fprintln(o.w, msg)
			logger.Info(msg)
		default: // TODO: cover. Never happen?
			msg := fmt.Sprintf("Cannot get.\tkey:%s error:%v", key, err)
			fmt.Fprintln(o.w, msg)
			logger.Info(msg)
		}
		return 1
	}

	msg := fmt.Sprintf("get\tkey:%s", key)
	logger.Info(msg)
	fmt.Fprintln(o.w, value)
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
			fmt.Fprintln(o.w, msg)
			logger.Info(msg)
		default: // TODO: cover. Never happen?
			msg := fmt.Sprintf("Cannot delete.\tkey:%s error:%v", key, err)
			fmt.Fprintln(o.w, msg)
			logger.Info(msg)
		}
		return 1 // TODO: cover
	}

	msg := fmt.Sprintf("delete\tkey:%s", key)
	logger.Info(msg)
	fmt.Fprintf(o.w, "%s is deleted.\n", key)
	return 0
}
