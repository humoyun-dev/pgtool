package cmd

import (
	"fmt"

	"github.com/humoyun-dev/pgcli/internal/sys"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "PostgreSQL ni o'rnatish va servisni ishga tushirish",
	RunE: func(cmd *cobra.Command, args []string) error {
		osName := sys.DetectOS()
		fmt.Println("Aniqlangan OS:", osName)

		switch osName {
		case "mac":
			if !sys.CommandExists("brew") {
				return fmt.Errorf("Homebrew topilmadi, avval uni o'rnat")
			}
			if err := sys.RunCmd("brew", "install", "postgresql@16"); err != nil {
				return err
			}
			return sys.RunCmd("brew", "services", "start", "postgresql@16")

		case "debian":
			if err := sys.RunCmd("sudo", "apt", "update"); err != nil {
				return err
			}
			if err := sys.RunCmd("sudo", "apt", "install", "-y", "postgresql", "postgresql-contrib"); err != nil {
				return err
			}
			return sys.RunCmd("sudo", "systemctl", "enable", "--now", "postgresql")

		default:
			return fmt.Errorf("bu OS uchun avtomatik install yo'q, qo'lda o'rnat")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
