package main

import (
	"os"

	"github.com/AbdessamadEnabih/Vertex/pkg/cmd"
)

func main() {

	root := cmd.NewRootCommand()

	serve := cmd.NewServeCommand()
	root.AddCommand(serve)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
