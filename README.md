# manssh

[![Build Status](https://travis-ci.org/xwjdsh/manssh.svg?branch=master)](https://travis-ci.org/xwjdsh/manssh)
[![Go Report Card](https://goreportcard.com/badge/github.com/xwjdsh/manssh)](https://goreportcard.com/report/github.com/xwjdsh/manssh)
[![codebeat badge](https://codebeat.co/badges/38954713-7443-4149-915d-4543da2a5da5)](https://codebeat.co/projects/github-com-xwjdsh-manssh-master)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)



manssh is a command line tool for manage your ssh alias config easily, inspire by [storm](https://github.com/emre/storm) project, power by Go.

Note: only `Host`, `HostName`, `User`, `Port`, `IdentityFile` and `ProxyCommand` is implemented at this point, **other configs and comments will loss**.

![](https://raw.githubusercontent.com/xwjdsh/manssh/master/screenshot/manssh.png)

## Install

#### Gopher
```shell
go get github.com/xwjdsh/manssh
```

#### Manual
Download it from [releases](https://github.com/xwjdsh/manssh/releases), and extact it to your `PATH` environment.

## Usage
```text
% manssh
NAME:
   manssh - Manage your ssh alias configs easily

USAGE:
   manssh [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
     add, a     add a new ssh alias record
     list, l    list all ssh alias records
     update, u  update existing ssh alias record
     delete, d  delete existing ssh alias record
     backup, b  backup ssh alias config records
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  (default: "/Users/wendell/.ssh/config")
   --help, -h                show help
   --version, -v             print the version

```

## Licence
[MIT License](https://github.com/xwjdsh/manssh/blob/master/LICENSE)
