package cmd

import (
	"fmt"

	"github.com/humoyun-dev/pgtool/internal/pg"
	"github.com/spf13/cobra"
)

var (
	createDbName  string
	createDbOwner string
)

var createDbCmd = &cobra.Command{
	Use:   "create-db",
	Short: "Yangi database yaratish",
	RunE: func(cmd *cobra.Command, args []string) error {
		if createDbName == "" || createDbOwner == "" {
			return fmt.Errorf("db nomi va owner kerak: --name --owner")
		}
		return pg.CreateDB(createDbName, createDbOwner)
	},
}

func init() {
	rootCmd.AddCommand(createDbCmd)

	createDbCmd.Flags().StringVar(&createDbName, "name", "", "Database name")
	createDbCmd.Flags().StringVar(&createDbOwner, "owner", "", "Owner user name")
}
