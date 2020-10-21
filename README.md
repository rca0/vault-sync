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

## Using homebrew to download binary

```bash
brew install rca0/tap/vault-sync
```

## Secret paths

Currently, this CLI tool support these vault secret paths

- kv-v2/
- secrets/

![image](https://user-images.githubusercontent.com/38728338/96727353-53a98d00-1389-11eb-9ce0-0ba5aa9e08b0.png)

## Params

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

## Running

```bash
vault-sync --srcaddr https://vault.domain --dstaddr https://vault-2.domain
```
