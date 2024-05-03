# Package
`github.com/tappoy/vault-cli`

# About
This vault cli command provides a simple way to interact with `github.com/tappoy/vault` package.

See [Usage](Usage.txt) for more details.

# Installation
```bash
go install github.com/tappoy/vault-cli@latest
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
