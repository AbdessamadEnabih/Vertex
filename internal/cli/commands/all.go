package commands

import (
	"fmt"

	"github.com/AbdessamadEnabih/Vertex/pkg/state"
	"github.com/spf13/cobra"
)

func All(globaleState *state.State) *cobra.Command {
	return &cobra.Command{
		Use:       "set",
		Short:     "Set a key-value pair",
		Example:   `set key value`,
		ValidArgs: []string{"key", "value"},
		Args:      cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			values := globaleState.GetAll()
			fmt.Println("All keys:", values)
		},
	}
}
