package commands

import (
	"fmt"

	"github.com/AbdessamadEnabih/Vertex/pkg/state"
	"github.com/spf13/cobra"
)

func NewGetCmd(globaleState *state.State) *cobra.Command {
	return &cobra.Command{
		Use:       "get",
		Short:     "Get a key-value pair",
		Example:   `get key`,
		ValidArgs: []string{"key"},
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value, err := globaleState.Get(key)
			if err != nil {
				fmt.Printf("Unable to get the key %v: %v\n", args[0], err)
			} else {
				fmt.Println("Value:", value)
			}
		},
	}
}
