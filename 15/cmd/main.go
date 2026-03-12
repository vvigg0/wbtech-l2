package main

import (
	"os"

	myshell "github.com/vvigg0/wbtech-l2/15/internal/shell"
)

func main() {
	sh := myshell.New(os.Stdin, os.Stdout, os.Stderr)
	code := sh.Run()
	os.Exit(code)
}
