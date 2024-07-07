// gomake setups your go project.
// go + make = gomake.
// Author: LitFill <marrazzy@gmail.com>
package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	flag "github.com/spf13/pflag"
	// "flag"

	"github.com/LitFill/fatal"
	"github.com/LitFill/gomake/templat"
)

type Config struct {
	name   string
	lib    bool
	logger bool
}

func usageMsg() {
	fmt.Println("usage: gomake --name Author/program\nflags:")
	flag.PrintDefaults()
}

func buatDirDanMasuk(nama string, loggr *slog.Logger) {
	fatal.Log(os.Mkdir(nama, 0o755), loggr,
		"Tidak dapat membuat directory", "dir", nama,
	)
	fatal.Log(os.Chdir(nama), loggr,
		"Tidak dapat pindah directory", "dir", nama,
	)
}

func validateName(logger *slog.Logger) bool {
	isValid := true
	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == "name" {
			name := f.Value.String()
			if name == "" {
				usageMsg()
				isValid = false
			}
			if len(strings.Split(name, "/")) < 2 {
				logger.Error("invalid name",
					"name", name,
					"correct format", "Author/program",
				)
				isValid = false
			}
		}
	})
	return isValid
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

func initMap(progName string, isLog bool, isLib bool) (map[string]string, map[string]string) {
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
	return peta, petaLib
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
	cfg := Config{}
	flag.StringVarP(&cfg.name, "name", "n", "", "name of the project [author/program]")
	flag.BoolVarP(&cfg.lib, "lib", "l", false, "create a lib instead of program")
	flag.BoolVarP(&cfg.logger, "logger", "g", false, "add logging using LitFill/fatal")
	flag.CommandLine.SortFlags = false

	flag.Parse()

	////////////////////////////////////////////////////////////////////////////////
	//                          mengklon flag agar aman                           //
	////////////////////////////////////////////////////////////////////////////////
	isLib := cfg.lib
	moduleName := cfg.name
	isLog := cfg.logger

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
	isValid := validateName(loggr)
	if !isValid {
		return
	}
	names := strings.Split(moduleName, "/")
	progName := names[len(names)-1]
	authorName := names[len(names)-2]

	data := newMetadata(authorName, moduleName, progName)

	// membuat direktori projek dan masuk ke dalamnya, melaporkan error lewat loggr
	buatDirDanMasuk(progName, loggr)

	// membuat map templat untuk dieksekusi
	peta, petaLib := initMap(progName, isLog, isLib)

	////////////////////////////////////////////////////////////////////////////////
	//            mengeksekusi templat dengan data dan membuat filenya            //
	////////////////////////////////////////////////////////////////////////////////
	if isLib {
		for nama, templ := range petaLib {
			fatal.Log(buatFileDenganTemplateDanEksekusi(nama, templ, data),
				loggr, "Tidak dapat mengeksekusi template",
				"file template", nama,
				"isLib", isLib,
			)
		}
	} else {
		for nama, templ := range peta {
			fatal.Log(buatFileDenganTemplateDanEksekusi(nama, templ, data),
				loggr, "Tidak dapat mengeksekusi template",
				"file template", nama,
				"isLib", isLib,
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
