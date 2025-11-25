package cmd

import (
	"fmt"

	"github.com/humoyun-dev/pgcli/internal/sys"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "PostgreSQL ni o'chirish (data ketishi mumkin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		osName := sys.DetectOS()
		fmt.Println("Aniqlangan OS:", osName)

		switch osName {
		case "mac":
			_ = sys.RunCmd("brew", "services", "stop", "postgresql@16")
			if err := sys.RunCmd("brew", "uninstall", "postgresql@16"); err != nil {
				return err
			}
			fmt.Println("Kerak bo'lsa data dir ni ham qo'lda o'chir: /opt/homebrew/var/postgresql@16")
		case "debian":
			_ = sys.RunCmd("sudo", "systemctl", "stop", "postgresql")
			_ = sys.RunCmd("sudo", "apt", "purge", "-y", "postgresql*")
			_ = sys.RunCmd("sudo", "rm", "-rf", "/var/lib/postgresql", "/etc/postgresql")
		default:
			return fmt.Errorf("bu OS uchun avtomatik uninstall yo'q")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
