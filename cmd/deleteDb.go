package cmd

import (
	"fmt"

	"github.com/humoyun-dev/pgcli/internal/pg"
	"github.com/spf13/cobra"
)

var deleteDbName string

var deleteDbCmd = &cobra.Command{
	Use:   "delete-db",
	Short: "Database o'chirish",
	RunE: func(cmd *cobra.Command, args []string) error {
		if deleteDbName == "" {
			return fmt.Errorf("db nomi kerak: --name")
		}
		return pg.DeleteDB(deleteDbName)
	},
}

func init() {
	rootCmd.AddCommand(deleteDbCmd)
	deleteDbCmd.Flags().StringVar(&deleteDbName, "name", "", "Database name")
}
