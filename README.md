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

# Operation Example
1. Make group `applications` for vault data accessors.
```bash
sudo groupadd applications
```

2. Make user `vault` for vault data maintenance.
```bash
sudo useradd -m -g applications -s /bin/bash vault
```

3. Add the `vault` user to the `applications` group.
```bash
sudo usermod -aG applications vault
```

4. Add the `vault` user to the `syslog` group.
```bash
sudo usermod -aG syslog vault
```

5. Change the gorup of the default vault directory.
```bash
sudo chgrp applications /srv
```

6. Change mode of the default vault directory.
```bash
sudo chmod 775 /srv
```



# License
[GPL-3.0](LICENSE)

# Author
[tappoy](https://github.com/tappoy)
