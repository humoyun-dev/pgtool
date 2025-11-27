package cmd

import (
	"fmt"

	"github.com/humoyun-dev/pgtool/internal/pg"
	"github.com/humoyun-dev/pgtool/internal/ui"
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
		interactive := ui.IsTerminal()

		if interactive {
			fmt.Println("=== Create Database ===")
			fmt.Println()
		}

		if interactive {
			if createDbName == "" {
				v, err := ui.Prompt("Database name: ")
				if err != nil {
					return err
				}
				createDbName = v
			}
			if createDbOwner == "" {
				v, err := ui.Prompt("Owner username: ")
				if err != nil {
					return err
				}
				createDbOwner = v
			}
		} else if createDbName == "" || createDbOwner == "" {
			return fmt.Errorf("Error: database name and owner are required: --name --owner (or interactive input).")
		}

		if createDbName == "" || createDbOwner == "" {
			return fmt.Errorf("Error: database name and owner are required: --name --owner (or interactive input).")
		}
		return pg.CreateDB(createDbName, createDbOwner)
	},
}

func init() {
	rootCmd.AddCommand(createDbCmd)

	createDbCmd.Flags().StringVar(&createDbName, "name", "", "Database name")
	createDbCmd.Flags().StringVar(&createDbOwner, "owner", "", "Owner user name")
}
