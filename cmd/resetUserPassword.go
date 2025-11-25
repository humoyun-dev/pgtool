package cmd

import (
	"fmt"

	"github.com/humoyun-dev/pgcli/internal/pg"
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
		if resetUserName == "" || resetUserPass == "" {
			return fmt.Errorf("username va yangi password kerak: --username --password")
		}
		return pg.ResetUserPassword(resetUserName, resetUserPass)
	},
}

func init() {
	rootCmd.AddCommand(resetUserPasswordCmd)

	resetUserPasswordCmd.Flags().StringVar(&resetUserName, "username", "", "User name")
	resetUserPasswordCmd.Flags().StringVar(&resetUserPass, "password", "", "New password")
}
