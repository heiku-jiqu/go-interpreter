package main

import (
	"fmt"
	"os"

	"github.com/heiku-jiqu/go-interpreter/repl"
)

func main() {
	fmt.Print("Welcome to the REPL!\n")
	r := repl.Repl{}
	r.Start(os.Stdin, os.Stdout)
}
