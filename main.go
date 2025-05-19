package main

import (
	"github.com/dr8co/monke/repl"
	"os/user"
)

func main() {
	usr, err := user.Current()

	if err != nil {
		panic(err)
	}

	repl.Start(usr.Username)
}
