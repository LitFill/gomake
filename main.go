// gomake setups your go project.
// go + make = gomake.
// Author: LitFill <marrazzy@gmail.com>
package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/LitFill/fatal"
	"github.com/LitFill/gomake/templat"
)

type Config struct {
	name   string
	lib    bool
	logger bool
}

func buatDirDanMasuk(nama string, loggr *slog.Logger) {
	fatal.Log(os.Mkdir(nama, 0o755), loggr,
		"Tidak dapat membuat directory", "dir", nama,
	)
	fatal.Log(os.Chdir(nama), loggr,
		"Tidak dapat pindah directory", "dir", nama,
	)
}

type MetadataProyek struct {
	AuthorName string
	ModuleName string
	ProgName   string
}

func newMetadata(author string, module string, program string) MetadataProyek {
	return MetadataProyek{
		AuthorName: author,
		ModuleName: module,
		ProgName:   program,
	}
}

type CmdsList []*exec.Cmd

func (c *CmdsList) add(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	*c = append(*c, cmd)
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
	////////////////////////////////////////////////////////////////////////////////
	//           membuat config untuk variable flag dan meparsingnya              //
	////////////////////////////////////////////////////////////////////////////////
	config := Config{}
	flag.BoolVar(&config.lib, "lib", false, "for creating a lib instead of program")
	flag.BoolVar(&config.lib, "l", false, "shorthand for -lib")
	flag.StringVar(&config.name, "name", "LitFill/test", "name of the project")
	flag.StringVar(&config.name, "n", "LitFill/test", "shorthand for -name")
	flag.BoolVar(&config.logger, "logger", false, "add logging using LitFill/fatal")
	flag.BoolVar(&config.logger, "g", false, "shorthand for -logger")

	flag.Parse()

	////////////////////////////////////////////////////////////////////////////////
	//                          mengklon flag agar aman                           //
	////////////////////////////////////////////////////////////////////////////////
	isLib := config.lib
	moduleName := config.name
	isLog := config.logger

	////////////////////////////////////////////////////////////////////////////////
	//                            setup LitFill/fatal                             //
	////////////////////////////////////////////////////////////////////////////////
	filename := filepath.Join(fatal.Assign(os.UserHomeDir())(
		slog.Default(), "cannot access user home dir"), ".gomake.log.json",
	)
	log_file := fatal.CreateLogFile(filename)
	defer log_file.Close()
	loggr := fatal.CreateLogger(io.MultiWriter(log_file, os.Stderr), slog.LevelInfo)

	////////////////////////////////////////////////////////////////////////////////
	//                             mengolah data nama                             //
	////////////////////////////////////////////////////////////////////////////////
	names := strings.Split(moduleName, "/")
	if len(names) < 2 || len(os.Args) < 2 {
		fatal.Log(exec.Command("gomake", "-h").Run(),
			loggr, "cannot run help command",
		)
		os.Exit(1)
	}

	progName := names[len(names)-1]
	authorName := names[len(names)-2]

	data := newMetadata(authorName, moduleName, progName)

	// membuat direktori projek dan masuk ke dalamnya, melaporkan error lewat loggr
	buatDirDanMasuk(progName, loggr)

	////////////////////////////////////////////////////////////////////////////////
	//                   membuat map templat untuk dieksekusi                     //
	////////////////////////////////////////////////////////////////////////////////
	peta := map[string]string{
		"main.go":   templat.MainTempl,
		"Makefile":  templat.MakeTempl,
		"README.md": templat.ReadmeTempl,
	}

	petaLib := map[string]string{
		progName + ".go": templat.LibTempl,
		"Makefile":       templat.LibMake,
	}

	if isLog {
		if !isLib {
			peta["main.go"] = templat.MainTemplWithLog
		} else {
			petaLib[progName+".go"] = templat.LibTemplWithLog
		}
	}

	////////////////////////////////////////////////////////////////////////////////
	//            mengeksekusi templat dengan data dan membuat filenya            //
	////////////////////////////////////////////////////////////////////////////////
	if isLib {
		for nama, templ := range petaLib {
			fatal.Log(buatFileDenganTemplateDanEksekusi(nama, templ, data),
				loggr, "Tidak dapat mengeksekusi template",
				"file template", nama,
			)
		}
	} else {
		for nama, templ := range peta {
			fatal.Log(buatFileDenganTemplateDanEksekusi(nama, templ, data),
				loggr, "Tidak dapat mengeksekusi template",
				"file template", nama,
			)
		}
	}

	////////////////////////////////////////////////////////////////////////////////
	//      membuat dan menjalankan list command untuk menyiapkan projek go       //
	////////////////////////////////////////////////////////////////////////////////
	cmdslist := make(CmdsList, 0)
	com := exec.Command

	if isLib {
		cmdslist.add(com("go", "mod", "init", "github.com/"+moduleName))
	} else {
		cmdslist.add(com("go", "mod", "init", moduleName))
	}
	if isLog {
		cmdslist.add(com("go", "get", "github.com/LitFill/fatal@latest"))
	}
	cmdslist.add(com("git", "init"))
	cmdslist.add(com("git", "add", "."))
	cmdslist.add(com("git", "commit", "-m", "initial commit"))

	for _, cmd := range cmdslist {
		fatal.Log(cmd.Run(), loggr,
			"Tidak dapat menjalankan perintah",
			"perintah", cmd.String(),
		)
	}

	////////////////////////////////////////////////////////////////////////////////
	//        mencetak pesan untuk user setelah selesai menyiapkan proyek         //
	////////////////////////////////////////////////////////////////////////////////
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
