// GopherVault is a local-first CLI knowledge manager built in Go.
//
// 📚 Go concept: "package main" tells Go this file should be compiled into
// an executable binary rather than a reusable shared library.
// The entry point of every Go executable is func main() in package main.
package main

import "github.com/Uvaisbugh/gophervault/cmd"

func main() {
	cmd.Execute()
}
