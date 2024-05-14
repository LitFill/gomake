package templat

var MainTempl = `// {{.ProgName}}, {{.AuthorName}} <author at email dot com>
// program for...
package main

import "fmt"

func main() {
	fmt.Println("Hello from {{.ModuleName}}!")
}
`
