package main

import (
	"fmt"
	"log"

	// "github.com/AbdessamadEnabih/Vertex/cli"
	"github.com/AbdessamadEnabih/Vertex/internal/net"
	"github.com/AbdessamadEnabih/Vertex/internal/persistance"
)

func main() {

	GlobalState, err := persistance.Load()

	if err != nil {
		fmt.Print(err)
		return
	}

	server := net.NewServer(GlobalState)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	defer persistance.Save(GlobalState)
}
