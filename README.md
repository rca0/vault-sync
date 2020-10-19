# Description

Simple and straight-forward tool to synchronize data from vault.

## Authentication

Currently, there are 2 auth methods supported

- github (default)
- token (https://www.vaultproject.io/docs/concepts/auth#tokens)

### Github auth method

Your vault cluster must be configured to accept this option, read more about it [here](https://www.vaultproject.io/docs/auth/github)

## How to use

Firstly you must build the tool, using Makefile just execute:

```bash
make build
```

```bash
NAME:
   vault-sync - copy vault data

USAGE:
    [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --srcaddr value   Source Vault Address
   --srctoken value  Source Vault Token
   --dstaddr value   Destination Vault Address
   --dsttoken value  Destination Vault Token
   --method value    Define auth method (github/token) (default: "github")
   --help, -h        show help
```

## To Run

```bash
./vault-sync --srcaddr https://vault.domain --dstaddr https://vault-2.domain
```
