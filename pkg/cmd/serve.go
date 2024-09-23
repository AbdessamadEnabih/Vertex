package cmd

import (
	"fmt"

	"github.com/AbdessamadEnabih/Vertex/pkg/server"
	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the Vertex server",
		Run: func(cmd *cobra.Command, args []string) {
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				fmt.Printf("Error parsing port flag: %v\n", err)
				return
			}

			fmt.Println("Starting Vertex server...")
			server.StartServer(port)
		},
	}

	// Add the port flag with a default value of ":6480"
	cmd.Flags().StringP("port", "p", "6480", "Port to run the Vertex server on")

	return cmd
}
