package main

import (
	"fmt"
	"log"

	"github.com/AbdessamadEnabih/Vertex/internal/network"
	"github.com/AbdessamadEnabih/Vertex/internal/persistance"
)

func main() {

	GlobalDataStore, err := persistance.Load()

	if err != nil {
		fmt.Print(err)
		return
	}

	server := network.NewServer(GlobalDataStore)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	defer persistance.Save(GlobalDataStore)
}
