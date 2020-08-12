# rkh
A command-line tool for removing entries of your ~/.ssh/known_hosts file


## Setup

### Pre-requisites

- Golang installend and configured in you machine

### Getting the tool

To install the CLI tool simply run:

```
go get github.com/ericovis/rkh
```

### Checking

If everything is configured properly you should be able to run `rkh --help` and get something like this:

```bash
$ rkh --help
Usage: rkh [parameters] [hostA hostB hostC ... hostZ]

Parameters:
  -force
        Suppress confirmation dialog.

```

## Contibuting

1. Fork this repo
2. Make changes
3. Open pull request
