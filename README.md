# gomake

Cli to start your new Go project. It uses git and makefile.

## Usage

```console
gomake "LitFill/program"
```

it makes `program` dir, inits `go mod init LitFill/program`, makes `main.go`, inits git repo, and creates `Makefile`.

then you can use the Makefile like so:

```console
make        # build for linux
make win    # build for windows
make run    # run the program
make help   # show help
```

README.md dalam bahasa Indonesia : [README.md](./README_ID.md)
