package main

import (
	"fmt"
	"github.com/dr8co/monke/repl"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Println("Hello", usr.Username, "This is the Monkey programming language!")
	fmt.Println("Feel free to type in commands")

	repl.Start(os.Stdin, os.Stdout)
}
