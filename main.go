// goMake is a Go package that makes things
// go + make = gomake
// Author: LitFill <marrazzy@gmail.com>
package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/LitFill/gomake/templat"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
	// Level: slog.LevelDebug,
	AddSource: true,
}))

func mayFatal[T comparable](val T, err error) T {
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	return val
}

func wrapErr(err error, msg string) error {
	return fmt.Errorf("%s, error: %w", msg, err)
}

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

	mayFatal(0, wrapErr(os.Mkdir(progName, 0o755), "Tidak dapat membuat directory"))
	mayFatal(0, wrapErr(
		os.Chdir(progName),
		fmt.Sprintf("Tidak dapat pindah ke ./%s/", progName),
	))

	peta := map[string]string{
		"main.go":   templat.MainTempl,
		"Makefile":  templat.MakeTempl,
		"README.md": templat.ReadmeTempl,
	}

	data := MetadataProyek{
		AuthorName: authorName,
		ModuleName: moduleName,
		ProgName:   progName,
	}

	for nama, templ := range peta {
		mayFatal(0, wrapErr(
			buatFileDenganTemplateDanEksekusi(nama, templ, data),
			fmt.Sprintf("Tidak dapat mengeksekusi %s", nama),
		))
	}

	cmds := make(Cmds)
	com := exec.Command

	cmds.add("init", com("go", "mod", "init", moduleName))
	cmds.add("git", com("git", "init"))
	cmds.add("git add", com("git", "add", "."))
	cmds.add("commit", com("git", "commit", "-m", "initial commit"))

	for name, cmd := range cmds {
		mayFatal(0, wrapErr(
			cmd.Run(),
			fmt.Sprintf("Tidak dapat menjalankan perintah %s", name),
		))
	}

	// pesan terakhir
	fmt.Println()
	fmt.Printf("Proyek %s telah dibuat\n", moduleName)
	fmt.Printf("Silahkan pindah direktori dengan menjalankan 'cd %s'\n", progName)
	fmt.Println("Disarankan untuk menjalankan 'go mod tidy'")
	fmt.Println("Jalankan program dengan 'make run'")
	fmt.Println("Compile dengan menjalankan 'make'")
}
