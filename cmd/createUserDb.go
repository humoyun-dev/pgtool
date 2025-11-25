package cmd

import (
	"fmt"

	"github.com/humoyun-dev/pgtool/internal/pg"
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
		if cudUser == "" || cudPass == "" || cudDb == "" {
			return fmt.Errorf("username, password va db nomi kerak: --username --password --db")
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
