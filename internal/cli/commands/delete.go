package commands

import (
	"fmt"

	"github.com/AbdessamadEnabih/Vertex/pkg/datastore"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(globaledatastore *datastore.DataStore) *cobra.Command {
	return &cobra.Command{
		Use:       "delete",
		Short:     "Delete a key-value pair",
		Example:   `delete key`,
		ValidArgs: []string{"key"},
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := globaledatastore.Delete(args[0])
			if err != nil {
				fmt.Printf("Failed to delete key %s : %v\n", args[0], err)
			}
		},
	}
}
