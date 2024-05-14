// goMake is a Go package that makes things
// go + make = gomake
// Author: LitFill <marrazzy@gmail.com>
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

type MetadataProyek struct {
	AuthorName string
	ModuleName string
	ProgName   string
}

func buatFileDenganTemplateDanEksekusi(namaFile, templ string, data MetadataProyek) error {
	file, err := os.Create(namaFile)
	if err != nil {
		return err
	}
	defer file.Close()

	t, err := template.New(namaFile).Parse(templ)
	if err != nil {
		return err
	}

	err = t.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gomake <module_name>")
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

	peta := map[string]string{
		"main.go":   mainTempl,
		"Makefile":  makeTempl,
		"README.md": readmeTempl,
	}

	data := MetadataProyek{
		AuthorName: authorName,
		ModuleName: moduleName,
		ProgName:   progName,
	}

	for nama, templ := range peta {
		err = buatFileDenganTemplateDanEksekusi(nama, templ, data)
		if err != nil {
			fmt.Printf("Tidak dapat membuat %s, error: %s\n", nama, err)
			os.Exit(1)
		}
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
