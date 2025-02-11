package commands

import (
	"fmt"

	"github.com/AbdessamadEnabih/Vertex/pkg/datastore"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(globaleDataStore *datastore.DataStore) *cobra.Command {
	return &cobra.Command{
		Use:       "update",
		Short:     "Update a key-value pair",
		Example:   `update key value`,
		ValidArgs: []string{"key", "value"},
		Args:      cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value := args[1]
			err := globaleDataStore.Update(key, value)
			if err != nil {
				fmt.Printf("Unable to update the key %v: %v\n", args[0], err)
			}
		},
	}
}
