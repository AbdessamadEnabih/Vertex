package commands

import (
	"fmt"

	"github.com/AbdessamadEnabih/Vertex/pkg/datastore"
	"github.com/spf13/cobra"
)

func NewFlushCmd(globaledatastore *datastore.DataStore) *cobra.Command {
	return &cobra.Command{
		Use:     "flush",
		Short:   "Flush the entire datastore",
		Example: `flush`,
		Run: func(cmd *cobra.Command, args []string) {
			err := globaledatastore.FlushAll()
			if err != nil {
				fmt.Printf("Failed to flush data: %v\n", err)
			}
		},
	}
}
