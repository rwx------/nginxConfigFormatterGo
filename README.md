# nginxConfigFormatterGo

go for nginx configure file formatter, a better and elegant nginx config formatter.

Inspired by <https://github.com/1connect/nginx-config-formatter.git>

## Features

- predictable formatted result.
- Comments on a separate line.
- neighbouring empty lines are collapsed to one empty line.
- curly braces placement follows Java convention.
- all lines are indented in uniform manner, with the given number spaces level (default 4).
- whitespaces are collapsed, except in comments and quotation marks.
- fixed multi-lines decompose problem (fixed for python version of [nginx-config-formatter](https://github.com/1connect/nginx-config-formatter.git) ).
- fixed quotes (`"`, `'`) match problem (fixed for python version of [nginx-config-formatter](https://github.com/1connect/nginx-config-formatter.git) ).

## Build Requirements

go 1.14.4+ (or go 1.13.12+)

## Installation

### 1. go get

```shell
go get github.com/rwx------/nginxConfigFormatterGo

# It may be installed at this path
$HOME/go/bin/nginxConfigFormatterGo
```

### 2. go build

```shell
git clone https://github.com/rwx------/nginxConfigFormatterGo.git
cd nginxConfigFormatterGo
go build
```

### 3. prebuild binary releases

You can get download prebuild binary at [Release Page](https://github.com/rwx------/nginxConfigFormatterGo/releases).

```shell
# linux
wget https://github.com/rwx------/nginxConfigFormatterGo/releases/download/v1.0.0/nginxConfigFormatterGo_linux_amd64 -O /usr/local/bin/nginxConfigFormatterGo

chmod +x /usr/local/bin/nginxConfigFormatterGo

# mac  
wget https://github.com/rwx------/nginxConfigFormatterGo/releases/download/v1.0.0/nginxConfigFormatterGo_darwin_amd64  -O /usr/local/bin/nginxConfigFormatterGo

chmod +x /usr/local/bin/nginxConfigFormatterGo
```

## Usage

```code
NAME:
   nginxConfigFormatterGo - Nginx config file formatter

USAGE:
   ./nginxConfigFormatterGo [-s 2] [-c utf-8] [-b] [-v] [-t] <filelists>

DESCRIPTION:
   Nginx config file formatter

AUTHOR:
   github.com/rwx------

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --charset value, -c value  current supported charset: gbk, gb18030, windows-1252, utf-8 (default: "utf-8")
   --space value, -s value    blank spaces indentation (default: 4)
   --backup, -b               backup the original config file
   --verbose, -v              verbose mode
   --testing, -t              perform a test run with no changes made
   --help, -h                 show help
```
