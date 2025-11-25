# pgtool

CLI for installing and managing PostgreSQL users & databases.

> Automates: install, uninstall, create users, create databases, list, reset passwords – from a single command.

## Features

- `install` / `uninstall` PostgreSQL (macOS Homebrew, Debian/Ubuntu)
- Create users with permissions (e.g. `CREATEDB`, `SUPERUSER`)
- Create databases with owner
- Create user + database in one command
- List users (`\du`) and databases (`\l`)
- Reset user password
- Delete users and databases

## Install

### From source

```bash
git clone https://github.com/<your-username>/pgtool.git
cd pgtool
go mod tidy
go build -o pgtool
sudo mv pgtool /usr/local/bin/
