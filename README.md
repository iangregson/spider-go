# spider-go

A simple web crawler.

A copy of [spider-monzo](https://github.com/iangregon/spider-monzo) in Go.

# Getting it running

Assuming you've got a Go toolchain installed.

Build it with

```bash
make
```

Then run it like 

```bash
$ bin/spider-go
A simple web crawler written in go.

Usage:
  spider-go [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  crawl       Crawl a given URL
  help        Help about any command

Flags:
  -h, --help     help for spider-go
  -t, --toggle   Help message for toggle

Use "spider-go [command] --help" for more information about a command.
```

There are tests too

```bash
make test
```

