package cmd

import (
	"fmt"

	"github.com/humoyun-dev/pgtool/internal/pg"
	"github.com/humoyun-dev/pgtool/internal/ui"
	"github.com/spf13/cobra"
)

var (
	resetUserName string
	resetUserPass string
)

var resetUserPasswordCmd = &cobra.Command{
	Use:   "reset-user-password",
	Short: "User parolini yangilash",
	RunE: func(cmd *cobra.Command, args []string) error {
		interactive := ui.IsTerminal()

		if interactive {
			fmt.Println("=== Reset User Password ===")
			fmt.Println()
		}

		if interactive {
			if resetUserName == "" {
				v, err := ui.Prompt("Username: ")
				if err != nil {
					return err
				}
				resetUserName = v
			}
			if resetUserPass == "" {
				v, err := ui.Prompt("New password: ")
				if err != nil {
					return err
				}
				resetUserPass = v
			}
		} else if resetUserName == "" || resetUserPass == "" {
			return fmt.Errorf("Error: username and new password are required: --username --password (or interactive input).")
		}

		if resetUserName == "" || resetUserPass == "" {
			return fmt.Errorf("Error: username and new password are required: --username --password (or interactive input).")
		}
		return pg.ResetUserPassword(resetUserName, resetUserPass)
	},
}

func init() {
	rootCmd.AddCommand(resetUserPasswordCmd)

	resetUserPasswordCmd.Flags().StringVar(&resetUserName, "username", "", "User name")
	resetUserPasswordCmd.Flags().StringVar(&resetUserPass, "password", "", "New password")
}
