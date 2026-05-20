package main

import (
	"os"
)

func main() {
	os.Exit(newCLI().run(os.Args[1:]))
}
