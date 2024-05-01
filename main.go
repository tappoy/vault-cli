package main

import (
	"fmt"
	"github.com/tappoy/logger"
	"github.com/tappoy/vault"
	"os"
	"path/filepath"
	"strings"
)

type option struct {
	command      string
	name         string
	password     string
	vaultDirRoot string
	logger       *logger.Logger
}

// print usage
func usage() {
	fmt.Printf(`Usage:
$ %s command [args] [name]

The commands are:
  help                       Show this help
  init [name]                Initialize a new vault
  set <key> <value> [name]   Set a key-value pair
  get <key> [name]           Get a value by key
  info [name]                Show information of the vault

You must give a password through the prompt when init, set and get.

Arguments:
  name - The name of the vault. Default is "vault".
  password - The password of the vault. It must be 8 to 32 characters.

Environment variables:
  VAULT_DIR - The directory of the vault. Default is "/srv".
  VAULT_LOG_DIR - The directory of the log. Default is "/var/log".
  VAULT_NAME - The name of the vault. Default is "vault".
`, os.Args[0])
	os.Exit(0)
}

func getName(nameIndex int) string {
	if len(os.Args) > nameIndex {
		return os.Args[nameIndex]
	} else if name := os.Getenv("VAULT_NAME"); strings.TrimSpace(name) != "" {
		return strings.TrimSpace(name)
	} else {
		return "vault"
	}
}

func getVaultDirRoot() string {
	if dir := os.Getenv("VAULT_DIR"); dir != "" {
		return dir
	} else {
		return "/srv"
	}
}

func getLogDirRoot() string {
	if dir := os.Getenv("VAULT_LOG_DIR"); dir != "" {
		return dir
	} else {
		return "/var/log"
	}
}

func newOptions(command string) *option {
	var name string

	switch command {
	case "help":
		usage()
	case "init", "info":
		name = getName(2)
	case "set":
		name = getName(4)
	case "get":
		name = getName(3)
	}

	if name == "" {
		fmt.Printf("Argument error. Run %s help\n", os.Args[0])
		return nil
	}

	logDir := filepath.Join(getLogDirRoot(), name)
	logger, err := logger.NewLogger(logDir)
	if err != nil {
		fmt.Printf("Cannot create logger '%s' %v\n", logDir, err)
		return nil
	}

	return &option{
		command:      command,
		name:         name,
		password:     "",
		vaultDirRoot: getVaultDirRoot(),
		logger:       logger,
	}
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
	if o.command == "get" || o.command == "set" {
		return os.Args[2]
	} else {
		panic("Something wrong")
	}
}

// get value
func (o *option) getValue() string {
	if o.command == "set" {
		if len(os.Args) >= 4 {
			return os.Args[3]
		} else {
			return ""
		}
	} else {
		panic("Something wrong")
	}
	return os.Args[3]
}

func main() {
	// check the arguments minimum length
	args := os.Args
	if len(args) < 2 {
		usage()
		os.Exit(1)
	}

	// parse arguments
	o := newOptions(args[1])
	if o == nil {
		os.Exit(1)
	}

	// info command
	if o.command == "info" {
		fmt.Println("[Vault Info]")
		fmt.Println("  name:", o.name)
		fmt.Println("  data:", o.getVaultDir())
		fmt.Println("  log :", o.logger.GetLogDir())
		fmt.Println("  init:", vault.IsInitialized(o.getVaultDir()))
		os.Exit(0)
	}

	// get password
	if err := o.getPassword(); err != nil {
		msg := fmt.Sprintf("Cannot get password %v", err)
		fmt.Println(msg)
		o.logger.Info(msg)
		os.Exit(1)
	}

	// create vault
	v, err := vault.NewVault(o.password, o.getVaultDir())
	if err != nil {
		switch err {
		case vault.ErrInvalidPasswordLength, vault.ErrPasswordIncorrect:
			msg := fmt.Sprintf("Wrong password.")
			fmt.Println(msg)
			o.logger.Notice(msg)
		default:
			msg := fmt.Sprintf("Cannot open vault %v [%s]", err, o.getVaultDir())
			fmt.Println(msg)
			o.logger.Info(msg)
		}
		os.Exit(1)
	}

	switch o.command {
	case "init":
		o.init(v)
	case "set":
		o.set(v)
	case "get":
		o.get(v)
	}

}

func (o *option) init(v *vault.Vault) {
	v, err := vault.NewVault(o.password, o.getVaultDir())
	if err != nil {
		body := ""
		switch err {
		case vault.ErrInvalidPasswordLength:
			body = fmt.Sprintf("password length must be 8 to 32 characters (length:%d)", len(o.password))
		default:
			body = fmt.Sprintf("%v", err)
		}
		header := "Cannot create vault: "
		fmt.Println(header + body)
		o.logger.Notice(header + body)
		os.Exit(1)
	}

	err = v.Init()
	if err != nil {
		msg := fmt.Sprintf("Cannot init vault %v [%s]", err, o.getVaultDir())
		fmt.Println(msg)
		o.logger.Notice(msg)
		os.Exit(1)
	}

	msg := fmt.Sprintf("Init vault %s", o.getVaultDir())
	fmt.Println(msg)
	o.logger.Notice(msg)
}

func (o *option) set(v *vault.Vault) {
	key := o.getKey()
	value := o.getValue()

	// check if the vault is initialized
	if !v.IsInitialized() {
		msg := fmt.Sprintf("Vault is not initialized [%s]", o.getVaultDir())
		fmt.Println(msg)
		o.logger.Info(msg)
		os.Exit(1)
	}

	if err := v.Set(key, value); err != nil {
		msg := fmt.Sprintf("Cannot set key '%s' '%v'", key, err)
		fmt.Println(msg)
		o.logger.Info(msg)
		os.Exit(1)
	}

	msg := fmt.Sprintf("Set key %s", key)
	fmt.Println(msg)
	o.logger.Info(msg)
}

func (o *option) get(v *vault.Vault) {
	key := o.getKey()
	value, err := v.Get(key)

	// check if the vault is initialized
	if !v.IsInitialized() {
		msg := fmt.Sprintf("Vault is not initialized [%s]", o.getVaultDir())
		fmt.Println(msg)
		o.logger.Info(msg)
		os.Exit(1)
	}

	if err != nil {
		switch err {
		case vault.ErrVariableNotFound:
			msg := fmt.Sprintf("Key '%s' not found", key)
			fmt.Println(msg)
			o.logger.Info(msg)
		default:
			msg := fmt.Sprintf("Cannot get key '%s' '%v'", key, err)
			fmt.Println(msg)
			o.logger.Info(msg)
		}
		os.Exit(1)
	}

	msg := fmt.Sprintf("Get key %s", key)
	o.logger.Info(msg)

	fmt.Println(value)
}
