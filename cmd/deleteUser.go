package cmd

import (
	"fmt"

	"github.com/humoyun-dev/pgtool/internal/pg"
	"github.com/spf13/cobra"
)

var deleteUserName string

var deleteUserCmd = &cobra.Command{
	Use:   "delete-user",
	Short: "Userni o'chirish",
	RunE: func(cmd *cobra.Command, args []string) error {
		if deleteUserName == "" {
			return fmt.Errorf("username kerak: --username")
		}
		return pg.DeleteUser(deleteUserName)
	},
}

func init() {
	rootCmd.AddCommand(deleteUserCmd)
	deleteUserCmd.Flags().StringVar(&deleteUserName, "username", "", "User name")
}
