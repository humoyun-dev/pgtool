package cmd

import (
	"fmt"

	"github.com/humoyun-dev/pgtool/internal/pg"
	"github.com/humoyun-dev/pgtool/internal/ui"
	"github.com/spf13/cobra"
)

var (
	cudUser string
	cudPass string
	cudPerm string
	cudDb   string
)

var createUserDbCmd = &cobra.Command{
	Use:   "create-user-db",
	Short: "User + DB ni bitta komandada yaratish",
	RunE: func(cmd *cobra.Command, args []string) error {
		interactive := ui.IsTerminal()

		if interactive {
			fmt.Println("=== Create User + Database ===")
			fmt.Println()
		}

		if interactive {
			if cudUser == "" {
				v, err := ui.Prompt("Username: ")
				if err != nil {
					return err
				}
				cudUser = v
			}
			if cudPass == "" {
				v, err := ui.Prompt("Password: ")
				if err != nil {
					return err
				}
				cudPass = v
			}
			if cudDb == "" {
				v, err := ui.Prompt("Database name: ")
				if err != nil {
					return err
				}
				cudDb = v
			}
		} else if cudUser == "" || cudPass == "" || cudDb == "" {
			return fmt.Errorf("Error: username, password and database name are required: --username --password --db (or interactive input).")
		}

		perms, err := choosePermissions(cudPerm, interactive)
		if err != nil {
			return err
		}
		cudPerm = perms

		if cudUser == "" || cudPass == "" || cudDb == "" {
			return fmt.Errorf("Error: username, password and database name are required: --username --password --db (or interactive input).")
		}
		return pg.CreateUserAndDB(cudUser, cudPass, cudPerm, cudDb)
	},
}

func init() {
	rootCmd.AddCommand(createUserDbCmd)

	createUserDbCmd.Flags().StringVar(&cudUser, "username", "", "User name")
	createUserDbCmd.Flags().StringVar(&cudPass, "password", "", "User password")
	createUserDbCmd.Flags().StringVar(&cudPerm, "perms", "", "Permissions (user uchun)")
	createUserDbCmd.Flags().StringVar(&cudDb, "db", "", "Database name")
}
