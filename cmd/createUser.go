package cmd

import (
	"fmt"

	"github.com/humoyun-dev/pgtool/internal/pg"

	"github.com/spf13/cobra"
)

var (
	createUserName string
	createUserPass string
	createUserPerm string
)

var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "PostgreSQL user yaratish",
	RunE: func(cmd *cobra.Command, args []string) error {
		if createUserName == "" || createUserPass == "" {
			return fmt.Errorf("username va password kerak: --username --password")
		}
		return pg.CreateUser(createUserName, createUserPass, createUserPerm)
	},
}

func init() {
	rootCmd.AddCommand(createUserCmd)

	createUserCmd.Flags().StringVar(&createUserName, "username", "", "User name")
	createUserCmd.Flags().StringVar(&createUserPass, "password", "", "Password")
	createUserCmd.Flags().StringVar(&createUserPerm, "perms", "", "Permissions (masalan: SUPERUSER, CREATEDB)")
}
