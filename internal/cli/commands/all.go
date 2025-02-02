package commands

import (
	"fmt"

	"github.com/AbdessamadEnabih/Vertex/pkg/datastore"
	"github.com/spf13/cobra"
)

func NewGetAllCmd(globaleDataStore *datastore.DataStore) *cobra.Command {
	return &cobra.Command{
		Use: 	 "all",
		Short:   "Get all keys",
		Example: `all`,
		Run: func(cmd *cobra.Command, args []string) {
			values := globaleDataStore.GetAll()
			fmt.Println("All keys:", values)
		},
	}
}
