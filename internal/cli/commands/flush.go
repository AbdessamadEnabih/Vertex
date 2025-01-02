package commands

import (
	"fmt"

	"github.com/AbdessamadEnabih/Vertex/pkg/state"
	"github.com/spf13/cobra"
)

func Flush(globaleState *state.State) *cobra.Command {
	return &cobra.Command{
		Use:     "flush",
		Short:   "Flush the entire state",
		Example: `flush`,
		Run: func(cmd *cobra.Command, args []string) {
			err := globaleState.FlushAll()
			if err != nil {
				fmt.Printf("Failed to flush data: %v\n", err)
			}
		},
	}
}
