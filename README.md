# manssh

[![Release](https://img.shields.io/github/release/xwjdsh/manssh.svg?style=flat-square)](https://github.com/xwjdsh/manssh/releases/latest)
[![Build Status](https://travis-ci.org/xwjdsh/manssh.svg?branch=master)](https://travis-ci.org/xwjdsh/manssh)
[![Go Report Card](https://goreportcard.com/badge/github.com/xwjdsh/manssh)](https://goreportcard.com/report/github.com/xwjdsh/manssh)
[![GoCover.io](https://img.shields.io/badge/gocover.io-89.0%25-green.svg)](https://gocover.io/github.com/xwjdsh/manssh)
[![GoDoc](https://godoc.org/github.com/xwjdsh/manssh?status.svg)](https://godoc.org/github.com/xwjdsh/manssh)
[![DUB](https://img.shields.io/dub/l/vibe-d.svg)](https://github.com/xwjdsh/manssh/blob/master/LICENSE)

manssh is a command line tool for managing your ssh alias config easily, inspired by [storm](https://github.com/emre/storm) project, powered by Go.

Note:<br/>
This project is actually a simple glue project, the most complex and core parsing ssh config file logic implements by [ssh_config](https://github.com/kevinburke/ssh_config), I didn't do much.<br/>
At first it was just a imitation of [storm](https://github.com/emre/storm), now it has become a little different.

![](https://raw.githubusercontent.com/xwjdsh/manssh/master/screenshot/manssh.gif)

## Feature

* No dependence.
* Add, list, query, delete ssh alias record.
* Backup ssh config.
* [Support Include directive.](#for-include-directive)

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

#### Docker
```shell
alias manssh='docker run -t --rm -v ~/.ssh/config:/root/.ssh/config wendellsun/manssh'
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
     backup, b  Backup SSH config files
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
# manssh list Test -ic
% manssh list test1 77
✔ Listing 1 records.

        test1 -> root@1.1.1.1:77
                identityfile = /Users/wendell/.ssh/wendell
```
It will display all alias records If no params offered, or it will using params as keywords query alias records.<br/>
If there is a `-it` option, it will ignore case when searching.

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

## For Include directive
If you use the `Include` directive, there are some extra notes.

Add `-p`(--path) flag for `list`,`add`,`update`,`delete` command to show the file path where the alias is located, it can also be set by the **MANSSH_SHOW_PATH** environment variable.

<details>
<summary><strong><code>MANSSH_SHOW_PATH</code></strong></summary>

Set to `true` to show the file path where the alias is located. Default is `false`.
</details>
<br/>

Add `-ap`(--addpath) flag for `add` command to specify the file path to which the alias is added, it can also be set by the **MANSSH_ADD_PATH** environment variable.

<details>
<summary><strong><code>MANSSH_ADD_PATH</code></strong></summary>

This file path indicates to which file to add the alias. Default is the entry config file.
</details>
<br/>

For convenience, you can export these environments in your `.zshrc` or `.bashrc`,
example:

```bash
export MANSSH_SHOW_PATH=true
export MANSSH_ADD_PATH=~/.ssh/config.d/temp
```

## Thanks
* [kevinburke/ssh_config](https://github.com/kevinburke/ssh_config)
* [urfave/cli](https://github.com/urfave/cli)
* [emre/storm](https://github.com/emre/storm)

## Licence
[MIT License](https://github.com/xwjdsh/manssh/blob/master/LICENSE)
