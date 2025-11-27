package cmd

import (
	"fmt"
	"strings"

	"github.com/humoyun-dev/pgtool/internal/pg"
	"github.com/humoyun-dev/pgtool/internal/ui"
	"github.com/spf13/cobra"
)

var (
	deleteUserName  string
	deleteUserForce bool
)

var deleteUserCmd = &cobra.Command{
	Use:   "delete-user",
	Short: "Userni o'chirish",
	RunE: func(cmd *cobra.Command, args []string) error {
		interactive := ui.IsTerminal()

		if interactive {
			fmt.Println("=== Delete User ===")
			fmt.Println()
		}

		if interactive {
			if deleteUserName == "" {
				roles, err := pg.ListRoleNames()
				if err != nil {
					return err
				}
				if len(roles) == 0 {
					return fmt.Errorf("Error: no users found to delete.")
				}
				selected, err := ui.SelectOne("Select user to delete", roles)
				if err != nil {
					return err
				}
				deleteUserName = selected
			}

			confirm, err := ui.Prompt(fmt.Sprintf("Are you sure you want to delete user '%s'? [y/N]: ", deleteUserName))
			if err != nil {
				return err
			}
			if !strings.HasPrefix(strings.ToLower(confirm), "y") {
				fmt.Println("Cancelled.")
				return nil
			}
		} else {
			if deleteUserName == "" || !deleteUserForce {
				return fmt.Errorf("Error: delete-user in non-interactive mode requires --username and --force.")
			}
		}
		return pg.DeleteUser(deleteUserName)
	},
}

func init() {
	rootCmd.AddCommand(deleteUserCmd)
	deleteUserCmd.Flags().StringVar(&deleteUserName, "username", "", "User name")
	deleteUserCmd.Flags().BoolVar(&deleteUserForce, "force", false, "Skip confirmation (non-interactive)")
}
