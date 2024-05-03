package main

import (
	"github.com/tappoy/logger"
	"github.com/tappoy/vault"
	ver "github.com/tappoy/version"

	"fmt"
	"io"
	"path/filepath"
)

type option struct {
	command      string
	name         string
	password     string
	vaultDirRoot string
	logger       *logger.Logger
	w            io.Writer
	args         []string
}

func (o *option) getVaultDir() string {
	return filepath.Join(o.vaultDirRoot, o.name)
}

// get password
func (o *option) getPassword() error {
	fmt.Print("Password: ")
	pwi := newPasswordInput()
	password, err := pwi.InputPassword()
	fmt.Print("\n")
	if err != nil {
		return err
	}
	o.password = password

	return nil
}

// get key
func (o *option) getKey() string {
	return o.args[2]
}

func (o *option) run() int {
	switch o.command {
	case "help":
		o.usage()
		return 0
	case "version":
		o.version()
		return 0
	case "genpw":
		o.generatePassword()
		return 0
	}

	if o.name == "" {
		fmt.Fprintf(o.w, "Argument error. Run %s help\n", o.args[0])
		return 1
	}

	// info command
	if o.command == "info" {
		fmt.Fprintln(o.w, "[Vault Info]")
		fmt.Fprintln(o.w, "  name:", o.name)
		fmt.Fprintln(o.w, "  data:", o.getVaultDir())
		fmt.Fprintln(o.w, "  log :", o.logger.GetLogDir())
		fmt.Fprintln(o.w, "  init:", vault.IsInitialized(o.getVaultDir()))
		return 0
	}

	// get password
	if err := o.getPassword(); err != nil {
		msg := fmt.Sprintf("Cannot get password.\terror:%v", err)
		fmt.Fprintln(o.w, msg)
		o.logger.Info(msg)
		return 1
	}

	// create vault
	v, err := vault.NewVault(o.password, o.getVaultDir())
	if err != nil {
		switch err {
		case vault.ErrInvalidPasswordLength, vault.ErrPasswordIncorrect:
			msg := fmt.Sprintf("Wrong password.")
			fmt.Fprintln(o.w, msg)
			o.logger.Notice(msg)
		default:
			msg := fmt.Sprintf("Cannot open vault.\terror:%v\tvaultDir:%s", err, o.getVaultDir())
			fmt.Fprintln(o.w, msg)
			o.logger.Info(msg)
		}
		return 1
	}

	switch o.command {
	case "init":
		return o.init(v)
	case "set":
		return o.set(v)
	case "get":
		return o.get(v)
	case "delete":
		return o.delete(v)
	default:
		return 1
	}
}

// print usage
func (o *option) usage() {
	fmt.Fprintf(o.w, `Usage:
$ vault-cli <command> [args...]

The commands are:
  help                       Show this help
  init [name]                Initialize a new vault
  set <key> <value> [name]   Set a key-value pair
  get <key> [name]           Get a value by key
  delete <key> [name]        Delete a key-value pair
  info [name]                Show information of the vault
  genpw                      Generate a random password
  version                    Show version

You must give a password through the prompt when init, set, get and delete.

Arguments:
  name - The name of the vault. Default is "vault".
  password - The password of the vault. It must be 8 to 32 characters.

Environment variables:
  VAULT_DIR - The root directory of the vault. Default is "/srv".
  VAULT_LOG_DIR - The root directory of the log. Default is "/var/log".
  VAULT_NAME - The name of the vault. Default is "vault". It will be used when the name argument is not given.

  By default, the vault data is stored in /srv/<name> and the log is stored in /var/log/<name>.
`)
}

// print version
func (o *option) version() {
	fmt.Fprintf(o.w, "vault-cli version %s\n", ver.Version())
}

// print random password
func (o *option) generatePassword() {
	password := vault.GeneratePassword()
	fmt.Fprintln(o.w, password)
}

func (o *option) init(v *vault.Vault) int {
	v, err := vault.NewVault(o.password, o.getVaultDir())
	if err != nil {
		body := ""
		switch err {
		case vault.ErrInvalidPasswordLength:
			body = fmt.Sprintf("password length must be 8 to 32 characters.\tlength:%d", len(o.password))
		default:
			body = fmt.Sprintf("error:%v", err)
		}
		header := "Cannot create vault. "
		fmt.Fprintln(o.w, header+body)
		o.logger.Notice(header + body)
		return 1
	}

	err = v.Init()
	if err != nil {
		msg := fmt.Sprintf("Cannot init vault.\terror:%v\tvaultDir:%s", err, o.getVaultDir())
		fmt.Fprintln(o.w, msg)
		o.logger.Notice(msg)
		return 1
	}

	msg := fmt.Sprintf("Init vault.\tvaultDir:%s", o.getVaultDir())
	fmt.Fprintln(o.w, msg)
	o.logger.Notice(msg)

	return 0
}

func (o *option) set(v *vault.Vault) int {
	key := o.getKey()

	var value string
	if len(o.args) >= 4 {
		value = o.args[3]
	} else {
		value = ""
	}

	// check if the vault is initialized
	if !v.IsInitialized() {
		msg := fmt.Sprintf("Vault is not initialized.\tvaultDir:%s", o.getVaultDir())
		fmt.Fprintln(o.w, msg)
		o.logger.Info(msg)
		return 1
	}

	if err := v.Set(key, value); err != nil {
		msg := fmt.Sprintf("Cannot set.\tkey:%s\terror:%v", key, err)
		fmt.Fprintln(o.w, msg)
		o.logger.Info(msg)
		return 1
	}

	msg := fmt.Sprintf("set\tkey:%s", key)
	o.logger.Info(msg)

	fmt.Fprintln(o.w, "Set successfully.")

	return 0
}

func (o *option) get(v *vault.Vault) int {
	key := o.getKey()
	value, err := v.Get(key)

	// check if the vault is initialized
	if !v.IsInitialized() {
		msg := fmt.Sprintf("Vault is not initialized.\tvaultDir:%s", o.getVaultDir())
		fmt.Fprintln(o.w, msg)
		o.logger.Info(msg)
		return 1
	}

	if err != nil {
		switch err {
		case vault.ErrKeyNotFound:
			msg := fmt.Sprintf("Not found.\tkey:%s", key)
			fmt.Fprintln(o.w, msg)
			o.logger.Info(msg)
		default:
			msg := fmt.Sprintf("Cannot get.\tkey:%s error:%v", key, err)
			fmt.Fprintln(o.w, msg)
			o.logger.Info(msg)
		}
		return 1
	}

	msg := fmt.Sprintf("get\tkey:%s", key)
	o.logger.Info(msg)

	fmt.Fprintln(o.w, value)
	return 0
}

func (o *option) delete(v *vault.Vault) int {
	key := o.getKey()

	// check if the vault is initialized
	if !v.IsInitialized() {
		msg := fmt.Sprintf("Vault is not initialized.\tvaultDir:%s", o.getVaultDir())
		fmt.Fprintln(o.w, msg)
		o.logger.Info(msg)
		return 1
	}

	if err := v.Delete(key); err != nil {
		msg := fmt.Sprintf("Cannot delete.\tkey:%s\terror:%v", key, err)
		fmt.Fprintln(o.w, msg)
		o.logger.Info(msg)
		return 1
	}

	msg := fmt.Sprintf("delete\tkey:%s", key)
	o.logger.Info(msg)

	fmt.Fprintf(o.w, "%s is deleted.\n", key)
	return 0
}
