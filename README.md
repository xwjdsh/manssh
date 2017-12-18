# manssh

[![Release](http://github-release-version.herokuapp.com/github/xwjdsh/manssh/release.svg?style=flat)](https://github.com/xwjdsh/manssh/releases/latest)
[![Build Status](https://travis-ci.org/xwjdsh/manssh.svg?branch=master)](https://travis-ci.org/xwjdsh/manssh)
[![Go Report Card](https://goreportcard.com/badge/github.com/xwjdsh/manssh)](https://goreportcard.com/report/github.com/xwjdsh/manssh)
[![GoCover.io](https://img.shields.io/badge/gocover.io-89.0%25-green.svg)](https://gocover.io/github.com/xwjdsh/manssh)
[![GoDoc](https://godoc.org/github.com/xwjdsh/manssh?status.svg)](https://godoc.org/github.com/xwjdsh/manssh)
[![DUB](https://img.shields.io/dub/l/vibe-d.svg)](https://github.com/xwjdsh/manssh/blob/master/LICENSE)

manssh is a command line tool for managing your ssh alias config easily, inspired by [storm](https://github.com/emre/storm) project, powered by Go.

![](https://raw.githubusercontent.com/xwjdsh/manssh/master/screenshot/manssh-12-16.gif)

## Feature

* No dependence.
* Add, list, query, delete ssh alias record.
* Backup ssh config.


## Install

#### Gopher
```shell
go get -u github.com/xwjdsh/manssh/cmd/manssh
```

#### Homebrew
```shell
brew tap xwjdsh/tap
brew install xwjdsh/tap/manssh
```

#### Manual
Download it from [releases](https://github.com/xwjdsh/manssh/releases), and extract it to your `PATH` directory.

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
     update, u  Update SSH record by specifying alias name
     delete, d  Delete SSH records by specifying alias names
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
# manssh add test1 root@1.1.1.1:77 -c IdentityFile=~/.ssh/wendell
% manssh add test1 root@1.1.1.1:77 -i ~/.ssh/wendell
✔ alias[test1] added successfully.

        test1 -> root@1.1.1.1:77
                identityfile = /Users/wendell/.ssh/wendell
```
Username and port config is optional, the username is current login username and port is `22` by default.<br/>
Using `-c` to set more config options. For convenience, `-i xxx` can instead of `-c identityfile=xxx`.

### List or query alias
```shell
# manssh list
# manssh list "*"
% manssh list test1 77
✔ Listing 1 records.

        test1 -> root@1.1.1.1:77
                identityfile = /Users/wendell/.ssh/wendell
```
It will display all alias records If no params offered, or it will using params as keywords query alias records. 

### Update an alias
```shell
# manssh update test1 -r test2
# manssh update test1 root@1.1.1.1:22022
% manssh update test1 -i "" -r test3 -c hostname=3.3.3.3 -c port=22022
✔ alias[test3] updated successfully.

        test3 -> root@3.3.3.3:22022
```
Update an existing alias record, it will replace origin user, hostname, port config's if connected string param offered.<br/>
You can use `-c` to update single and extra config option, `-c identityfile= -c proxycommand=` will remove `identityfile` and `proxycommand` options. <br/>
For convenience, `-i xxx` can instead of `-c identityfile=xxx`<br/>
Rename the alias specified by `-r` flag.

### Delete one or more alias
```shell
# manssh delete test1
% manssh delete test1 test2
✔ alias[test1,test2] deleted successfully.
```

### Backup ssh config
```
% manssh backup ./config_backup
✔ backup ssh config to [./config_backup] successfully.
```

## Thanks
* [kevinburke/ssh_config](https://github.com/kevinburke/ssh_config)
* [urfave/cli](https://github.com/urfave/cli)
* [emre/storm](https://github.com/emre/storm)

## Licence
[MIT License](https://github.com/xwjdsh/manssh/blob/master/LICENSE)
