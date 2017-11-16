# manssh

[![Release](http://github-release-version.herokuapp.com/github/xwjdsh/manssh/release.svg?style=flat)](https://github.com/xwjdsh/manssh/releases/latest)
[![Build Status](https://travis-ci.org/xwjdsh/manssh.svg?branch=master)](https://travis-ci.org/xwjdsh/manssh)
[![Go Report Card](https://goreportcard.com/badge/github.com/xwjdsh/manssh)](https://goreportcard.com/report/github.com/xwjdsh/manssh)
[![GoCover.io](https://img.shields.io/badge/gocover.io-89.0%25-green.svg)](https://gocover.io/github.com/xwjdsh/manssh)
[![GoDoc](https://godoc.org/github.com/xwjdsh/manssh?status.svg)](https://godoc.org/github.com/xwjdsh/manssh)
[![DUB](https://img.shields.io/dub/l/vibe-d.svg)](https://github.com/xwjdsh/manssh/blob/master/LICENSE)

manssh is a command line tool for managing your ssh alias config easily, inspired by [storm](https://github.com/emre/storm) project, powered by Go.

![](https://raw.githubusercontent.com/xwjdsh/manssh/master/screenshot/manssh.gif)

## Feature

* No dependence.
* Add, list, query, delete ssh alias record.
* Backup ssh config.


## Install

#### Gopher
```shell
go get github.com/xwjdsh/manssh/cmd/manssh
```

#### Homebrew
```shell
brew tap xwjdsh/tap
brew install xwjdsh/tap/manssh
```

#### Manual
Download it from [releases](https://github.com/xwjdsh/manssh/releases), and extact it to your `PATH` directory.

## Usage
```text
% manssh
NAME:
   manssh - Manage your ssh alias configs easily

USAGE:
   manssh [global options] command [command options] [arguments...]

VERSION:
   master

COMMANDS:
     add, a     Add a new SSH alias record
     list, l    List or query SSH alias records
     update, u  Update SSH record by specified alias name
     delete, d  Delete SSH records by specified alias names
     backup, b  Backup SSH alias config records
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value  (default: "/Users/wendell/.ssh/config")
   --help, -h              show help
   --version, -v           print the version
```

### Add a new alias
```shell
# manssh add test2 2.2.2.2
% manssh add test1 root@1.1.1.1:77 -c IdentityFile=~/.ssh/wendell
 success alias[test1] added successfully.

        test1 -> root@1.1.1.1:77
                IdentityFile = ~/.ssh/wendell
```
Using `-c` to set more config options, username and port config is optional, the username is current login username and port is `22` by default.

### List or query alias
```shell
# manssh list
# manssh list "*"
% manssh list test1 77
 success Listing 1 records.

        test1 -> root@1.1.1.1:77
                IdentityFile = ~/.ssh/wendell
```
It will display all alias records If no params offered, or it will using params as keywords query alias records. 

### Update a alias
```shell
# manssh update test1 -r test2
# manssh update test1 root@1.1.1.1:22022
% manssh update test1 -r test3 -c user=wendell -c port=22022
 success alias[test3] updated successfully.

        test3 -> wendell@1.1.1.1:22022
                identityfile = ~/.ssh/wendell
```
Update existing alias record, it will replace origin user, hostname, port config if connect string param offered, you can using `-c` to update single and extra config option. Rename alias specified by `-r` flag.

### Delete one or more alias
```shell
# manssh delete test1
% manssh update test1 test2
 success alias[test1,test2] deleted successfully.
```

### Backup ssh config
```
% manssh backup ./config_backup
 success backup ssh config to [./config_backup] successfully.
```

## Licence
[MIT License](https://github.com/xwjdsh/manssh/blob/master/LICENSE)
