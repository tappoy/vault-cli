# Package
`github.com/tappoy/vault-cli`

# About
This vault cli command provides a simple way to interact with `github.com/tappoy/vault` package.

# Installation
```bash
go install github.com/tappoy/vault-cli@latest
```

# Usage
```
Usage:
$ vault-cli command [args] [name]

The commands are:
  help                       Show this help
  init [name]                Initialize a new vault
  set <key> <value> [name]   Set a key-value pair
  get <key> [name]           Get a value by key
  info [name]                Show information of the vault

You must give a password through the prompt when init, set and get.

Arguments:
  name - The name of the vault. Default is `vault`.
  password - The password of the vault. It must be 8 to 32 characters.

Environment variables:
  VAULT_DIR - The directory of the vault. Default is `/srv`.
  VAULT_LOG_DIR - The directory of the log. Default is `/var/log`.
  VAULT_NAME - The name of the vault. Default is `vault`.
```

# License
[GPL-3.0](LICENSE)

# Author
[tappoy](https://github.com/tappoy)
