# gomake

Cli to start your new Go project. It uses git and makefile.

## Installation

### go install 

use go to install `gomake` :

```bash
go install github.com/LitFill/gomake@latest
```

### dependencies

`gomake` uses git and makefile so it requires `git` and `make` to be installed.

## Usage

```bash
gomake "LitFill/program"
```

it makes `program` dir, inits `go mod init LitFill/program`, makes `main.go`, inits git repo, and creates `Makefile`.

then you can use the Makefile like so:

```bash
make        # build for linux
make win    # build for windows
make run    # run the program
make help   # show help
```

README.md dalam bahasa Indonesia : [README.md](./README_ID.md)
