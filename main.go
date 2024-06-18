// goMake is a Go package that makes things
// go + make = gomake
// Author: LitFill <marrazzy@gmail.com>
package main

import (
	"flag"
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
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s, error: %w", msg, err)
}

func fatalWrapf(err error, format string, a ...any) {
	mayFatal(0, wrapErr(err, fmt.Sprintf(format, a...)))
}

type CmdsList []*exec.Cmd

func (c *CmdsList) add(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	*c = append(*c, cmd)
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
	var (
		lib  bool
		name string
	)
	flag.BoolVar(&lib, "lib", false, "for creating a lib instead of program")
	flag.BoolVar(&lib, "l", false, "shorthand for -lib")
	flag.StringVar(&name, "name", "LitFill/test", "name of the project")
	flag.StringVar(&name, "n", "LitFill/test", "shorthand for -name")

	flag.Parse()

	isLib := lib
	moduleName := name
	names := strings.Split(moduleName, "/")

	if len(names) < 2 || len(os.Args) < 2 {
		m := `Usage: gomake -n <module_name> [flags]
module_name = 'author/program'

flag options:
	-lib  or -l:   create library instead
	-name or -n:   the name of the module in format of "author/program"`
		fmt.Println(m)
		os.Exit(1)
	}

	progName := names[len(names)-1]
	authorName := names[len(names)-2]

	fatalWrapf(os.Mkdir(progName, 0o755), "Tidak dapat membuat directory %s", progName)
	fatalWrapf(os.Chdir(progName), "Tidak dapat pindah ke ./%s/", progName)

	peta := map[string]string{
		"main.go":   templat.MainTempl,
		"Makefile":  templat.MakeTempl,
		"README.md": templat.ReadmeTempl,
	}

	petaLib := map[string]string{
		fmt.Sprintf("%s.go", progName): templat.LibTempl,
		"Makefile":                     templat.LibMake,
	}

	data := MetadataProyek{
		AuthorName: authorName,
		ModuleName: moduleName,
		ProgName:   progName,
	}

	if isLib {
		for nama, templ := range petaLib {
			fatalWrapf(buatFileDenganTemplateDanEksekusi(nama, templ, data),
				"Tidak dapat mengeksekusi %s", nama,
			)
		}
	} else {
		for nama, templ := range peta {
			fatalWrapf(buatFileDenganTemplateDanEksekusi(nama, templ, data),
				"Tidak dapat mengeksekusi %s", nama,
			)
		}
	}

	cmdslist := make(CmdsList, 0)
	com := exec.Command

	if isLib {
		cmdslist.add(com("go", "mod", "init", fmt.Sprintf("github.com/%s", moduleName)))
	} else {
		cmdslist.add(com("go", "mod", "init", moduleName))
	}
	cmdslist.add(com("git", "init"))
	cmdslist.add(com("git", "add", "."))
	cmdslist.add(com("git", "commit", "-m", "initial commit"))

	for _, cmd := range cmdslist {
		fatalWrapf(cmd.Run(), "Tidak dapat menjalankan perintah %s", cmd.String())
	}

	// pesan terakhir
	fmt.Println()
	if isLib {
		fmt.Printf("-- Proyek github.com/%s telah dibuat\n", moduleName)
	} else {
		fmt.Printf("-- Proyek %s telah dibuat\n", moduleName)
	}
	fmt.Printf("-- Silahkan pindah direktori dengan menjalankan 'cd %s'\n", progName)
	fmt.Println("-- Disarankan untuk menjalankan 'go mod tidy'")
	if !isLib {
		fmt.Println("-- Jalankan program dengan 'make run'")
		fmt.Println("-- Compile dengan menjalankan 'make'")
	}
}
