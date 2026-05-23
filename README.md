# baum

> *Baum* is German for "tree."

A rich terminal tree viewer written in Go, inspired by [lstr](https://github.com/bgreenwell/lstr).

## Install

```bash
go install github.com/go-spass/baum@latest
```

## Usage

```bash
baum [path]               # tree view of current or given directory
baum -L 2                 # limit depth to 2 levels
baum -a                   # include hidden files
baum -d                   # directories only
baum --color never        # disable color output
```

## Features

- [ ] Classic tree view with Unicode branch characters
- [ ] Depth limiting (`-L`)
- [ ] Hidden file toggle (`-a`)
- [ ] File sizes, permissions, icons
- [ ] Git status integration (`-G`)
- [ ] `.gitignore` awareness (`-g`)
- [ ] Flexible sorting
- [ ] Interactive TUI mode

## Development

```bash
make build   # build ./baum binary
make test    # run tests
make install # install to $GOPATH/bin
```
