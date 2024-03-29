/*
goMake is a Go package that makes things
Author: LitFill <marrazzy@gmail.com>
*/
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type Cmds map[string]*exec.Cmd

func (c *Cmds) add(name string, cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	(*c)[name] = cmd
}

var mainTempl = `/*
  {{.ProgName}} {{.AuthorName}} <email>
*/
package main

import "fmt"

func main() {
	fmt.Println("Hello from {{.ModuleName}}!")
}
`

var makeTempl = `COMPILER := go
BINNAME := {{.ProgName}}

BUILDCMD := $(COMPILER) build
OUTPUT := -o $(BINNAME)
FLAGS := -v

RUNCMD := $(COMPILER) run

.PHONY: all build run clean win help

all: build ## Build the binary for Linux

build: main.go ## Actually build the binary
	@echo "Building $(BINNAME) for Linux"
	@$(BUILDCMD) $(OUTPUT) $(FLAGS)

win: main.go ## Build the binary for Windows
	@echo "Building $(BINNAME) for Windows"
	@$(BUILDCMD) $(OUTPUT).exe $(FLAGS)

run: main.go ## Run the main.go
	@echo "Running $(BINNAME)"
	@$(RUNCMD) $(FLAGS) $^

clean: ## Clean up
	@echo "Cleaning up"
	@rm -f $(BINNAME)*

help: ## Prints help for targets with comments
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: goMake <module_name>")
		os.Exit(1)
	}

	moduleName := os.Args[1]
	names := strings.Split(moduleName, "/")
	progName := names[len(names)-1]
	authorName := names[len(names)-2]

	err := os.MkdirAll(progName, 0o755)
	if err != nil {
		fmt.Println("Tidak dapat membuat directory, error:", err)
		os.Exit(1)
	}

	err = os.Chdir(progName)
	if err != nil {
		fmt.Printf("Tidak dapat pindah ke ./%s, error: %s\n", progName, err)
		os.Exit(1)
	}

	mainFile, err := os.Create("main.go")
	if err != nil {
		fmt.Println("Tidak bisa membuat file main.go, error:", err)
		os.Exit(1)
	}
	defer mainFile.Close()

	makeFile, err := os.Create("Makefile")
	if err != nil {
		fmt.Println("Tidak bisa membuat file Makefile, error:", err)
		os.Exit(1)
	}
	defer makeFile.Close()

	templMain, err := template.New("main").Parse(mainTempl)
	if err != nil {
		fmt.Println("Tidak dapat membuat template, error:", err)
		os.Exit(1)
	}

	templMake, err := template.New("make").Parse(makeTempl)
	if err != nil {
		fmt.Println("Tidak dapat membuat template, error:", err)
		os.Exit(1)
	}

	data := struct {
		AuthorName string
		ModuleName string
		ProgName   string
	}{
		AuthorName: authorName,
		ModuleName: moduleName,
		ProgName:   progName,
	}

	err = templMain.Execute(mainFile, data)
	if err != nil {
		fmt.Println("Tidak dapat membuat template, error:", err)
		os.Exit(1)
	}

	err = templMake.Execute(makeFile, data)
	if err != nil {
		fmt.Println("Tidak dapat membuat template, error:", err)
		os.Exit(1)
	}

	cmds := make(Cmds)
	com := exec.Command

	cmds.add("init", com("go", "mod", "init", moduleName))
	cmds.add("git", com("git", "init"))
	cmds.add("git add", com("git", "add", "."))
	cmds.add("commit", com("git", "commit", "-m", "initial commit"))

	for name, cmd := range cmds {
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Tidak dapat menjalankan perintah %s, error: %s\n", name, err)
			os.Exit(1)
		}
	}

	// pesan terakhir

	fmt.Println()
	fmt.Printf("Proyek %s telah dibuat\n", moduleName)
	fmt.Printf("Silahkan pindah direktori dengan menjalankan 'cd %s'\n", progName)
	fmt.Println("Disarankan untuk menjalankan 'go mod tidy'")
	fmt.Println("Jalankan program dengan 'make run'")
	fmt.Println("Compile dengan menjalankan 'make'")
}
