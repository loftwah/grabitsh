package main

import (
	"fmt"
	"os"

	"github.com/loftwah/grabitsh/cmd/grabitsh"
)

func main() {
	if err := grabitsh.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
