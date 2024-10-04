package main

import (
	"os"

	"github.com/AbdessamadEnabih/Vertex/pkg/cmd"
)

func main() {

	root := cmd.NewRootCommand()


	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
