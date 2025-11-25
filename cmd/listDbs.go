package cmd

import (
	"github.com/humoyun-dev/pgtool/internal/pg"
	"github.com/spf13/cobra"
)

var listDbsCmd = &cobra.Command{
	Use:   "list-dbs",
	Short: "Database larni ko'rsatish",
	RunE: func(cmd *cobra.Command, args []string) error {
		return pg.ListDBs()
	},
}

func init() {
	rootCmd.AddCommand(listDbsCmd)
}
