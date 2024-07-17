package templat

var MainTempl = `// {{ .ProgName }}, {{ .AuthorName }} <author at email dot com>
// program for ...
package main

import "fmt"

func main() {
	fmt.Println("Hello from {{ .ModuleName }}!")
}
`

var MainTemplWithLog = `// {{ .ProgName }}, {{ .AuthorName }} <author at email dot com>
// program for ...
package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/LitFill/fatal"
)

func main() {
	logFile := fatal.CreateLogFile("log.json")
	defer logFile.Close()
	logger := fatal.CreateLogger(io.MultiWriter(logFile,os.Stderr), slog.LevelInfo)
	logger.Info("Hello from {{ .ModuleName }}!")
}
`
