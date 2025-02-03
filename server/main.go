package main

import (
	"fmt"
	"log"

	"github.com/AbdessamadEnabih/Vertex/internal/network"
	"github.com/AbdessamadEnabih/Vertex/pkg/persistence"
)

func main() {

	GlobalDataStore, err := persistence.Load()

	if err != nil {
		fmt.Print(err)
		return
	}

	server := network.NewServer(GlobalDataStore)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	defer persistence.Save(GlobalDataStore)
}
