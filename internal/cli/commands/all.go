package commands

import (
	"fmt"

	"github.com/AbdessamadEnabih/Vertex/pkg/state"
	"github.com/spf13/cobra"
)

func NewGetAllCmd(globaleState *state.State) *cobra.Command {
	return &cobra.Command{
		Use: 	 "all",
		Short:   "Get all keys",
		Example: `all`,
		ValidArgs: []string{"key", "value"},
		Args:      cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			values := globaleState.GetAll()
			fmt.Println("All keys:", values)
		},
	}
}
