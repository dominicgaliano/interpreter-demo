package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dominicgaliano/interpreter-demo/repl"
)

func main() {
    user, err := user.Current()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Welcome %s, this is the Monkey programming language!\n", user.Username)
    fmt.Printf("Input commands below:\n")
    repl.Start(os.Stdin, os.Stdout)
}
