package cmd

import (
	"github.com/humoyun-dev/pgcli/internal/pg"
	"github.com/spf13/cobra"
)

var listUsersCmd = &cobra.Command{
	Use:   "list-users",
	Short: "User (role) larni ko'rsatish",
	RunE: func(cmd *cobra.Command, args []string) error {
		return pg.ListUsers()
	},
}

func init() {
	rootCmd.AddCommand(listUsersCmd)
}
