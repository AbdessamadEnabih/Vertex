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
			fmt.Println("Starting Vertex server...")
			server.StartServer()
		},
	}

	cmd.Flags().StringP("port", "p", "8080", "Specify the port to listen on")
	cmd.Flags().BoolP("debug", "d", false, "Enable debug mode")

	return cmd
}
