package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pgcli",
	Short: "PostgreSQL o'rnatish va boshqarish uchun CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// RootCommand testlar uchun root komandani qaytaradi.
func RootCommand() *cobra.Command {
	return rootCmd
}
