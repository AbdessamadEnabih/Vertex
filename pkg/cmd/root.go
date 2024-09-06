package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "vertex",
		Short: "open-source authorization service & policy engine based on Google Zanzibar.",
		Long:  "open-source authorization service & policy engine based on Google Zanzibar.",
	}
}
