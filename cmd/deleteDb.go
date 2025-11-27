package cmd

import (
	"fmt"
	"strings"

	"github.com/humoyun-dev/pgtool/internal/pg"
	"github.com/humoyun-dev/pgtool/internal/ui"
	"github.com/spf13/cobra"
)

var (
	deleteDbName  string
	deleteDbForce bool
)

var deleteDbCmd = &cobra.Command{
	Use:   "delete-db",
	Short: "Database o'chirish",
	RunE: func(cmd *cobra.Command, args []string) error {
		interactive := ui.IsTerminal()

		if interactive {
			fmt.Println("=== Delete Database ===")
			fmt.Println()
		}

		if interactive {
			if deleteDbName == "" {
				dbs, err := pg.ListDatabaseNames()
				if err != nil {
					return err
				}
				if len(dbs) == 0 {
					return fmt.Errorf("Error: no databases found to delete.")
				}
				selected, err := ui.SelectOne("Select database to delete", dbs)
				if err != nil {
					return err
				}
				deleteDbName = selected
			}

			confirm, err := ui.Prompt(fmt.Sprintf("Are you sure you want to delete database '%s'? [y/N]: ", deleteDbName))
			if err != nil {
				return err
			}
			if !strings.HasPrefix(strings.ToLower(confirm), "y") {
				fmt.Println("Cancelled.")
				return nil
			}
		} else {
			if deleteDbName == "" || !deleteDbForce {
				return fmt.Errorf("Error: delete-db in non-interactive mode requires --name and --force.")
			}
		}
		return pg.DeleteDB(deleteDbName)
	},
}

func init() {
	rootCmd.AddCommand(deleteDbCmd)
	deleteDbCmd.Flags().StringVar(&deleteDbName, "name", "", "Database name")
	deleteDbCmd.Flags().BoolVar(&deleteDbForce, "force", false, "Skip confirmation (non-interactive)")
}
