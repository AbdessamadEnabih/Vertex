package commands

import (
	"fmt"

	"github.com/AbdessamadEnabih/Vertex/pkg/datastore"
	"github.com/spf13/cobra"
)

func NewSetCmd(globaledatastore *datastore.DataStore) *cobra.Command {
	return &cobra.Command{
		Use:       "set",
		Short:     "Set a key-value pair",
		Example:   `set key value`,
		ValidArgs: []string{"key", "value"},
		Args:      cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value := args[1]
			err := globaledatastore.Set(key, value)
			if err != nil {
				fmt.Printf("Unable to set the key %v: %v\n", args[0], err)
			}
		},
	}
}
