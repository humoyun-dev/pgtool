package cmd

import (
	"fmt"
	"strings"

	"github.com/humoyun-dev/pgtool/internal/sys"
	"github.com/humoyun-dev/pgtool/internal/ui"
	"github.com/spf13/cobra"
)

var (
	uninstallHard bool
	uninstallYes  bool
)

type stepResult struct {
	label string
	err   error
}

func runIfExists(label, name string, args ...string) stepResult {
	if !sys.CommandExists(name) {
		return stepResult{label: label, err: fmt.Errorf("%s not found", name)}
	}
	return stepResult{label: label, err: sys.RunCommand(label, name, args...)}
}

func uninstallMac(hard bool) []stepResult {
	var results []stepResult
	results = append(results, runIfExists("stop postgresql@16 service", "brew", "services", "stop", "postgresql@16"))
	results = append(results, runIfExists("stop postgresql service", "brew", "services", "stop", "postgresql"))
	results = append(results, runIfExists("uninstall postgresql@16", "brew", "uninstall", "postgresql@16"))
	results = append(results, runIfExists("uninstall postgresql", "brew", "uninstall", "postgresql"))

	if hard {
		paths := []string{
			"/opt/homebrew/var/postgresql@16",
			"/usr/local/var/postgres",
			fmt.Sprintf("%s/Library/Caches/pgtool", sys.HomeDir()),
		}
		for _, p := range paths {
			label := "remove " + p
			results = append(results, stepResult{label: label, err: sys.RemovePath(label, p)})
		}
	}
	return results
}

func uninstallDebian(hard bool) []stepResult {
	var results []stepResult
	if sys.CommandExists("systemctl") {
		results = append(results, runIfExists("stop postgresql service", "sudo", "systemctl", "stop", "postgresql"))
	} else if sys.CommandExists("service") {
		results = append(results, runIfExists("stop postgresql service", "sudo", "service", "postgresql", "stop"))
	}
	if sys.CommandExists("apt-get") {
		results = append(results, runIfExists("remove packages", "sudo", "apt-get", "remove", "--purge", "postgresql*", "-y"))
	}

	if hard {
		paths := []string{
			"/var/lib/postgresql",
			"/var/log/postgresql",
			"/etc/postgresql",
		}
		for _, p := range paths {
			label := "remove " + p
			results = append(results, stepResult{label: label, err: sys.RemovePath(label, p)})
		}
	}
	return results
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "PostgreSQL ni o'chirish (data ketishi mumkin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if uninstallHard {
			if ui.IsTerminal() && !uninstallYes {
				fmt.Println("=== PostgreSQL HARD UNINSTALL ===")
				fmt.Println("This will stop PostgreSQL, uninstall it, and remove data directories, logs, and caches.")
				fmt.Println("This may delete ALL your PostgreSQL databases on this machine.")
				fmt.Println()
				ans, err := ui.Prompt(`Type "DELETE EVERYTHING" to continue: `)
				if err != nil {
					return err
				}
				if strings.TrimSpace(ans) != "DELETE EVERYTHING" {
					fmt.Println("Cancelled.")
					return nil
				}
			} else if !ui.IsTerminal() && !uninstallYes {
				return fmt.Errorf("Error: --hard uninstall in non-interactive mode requires --yes")
			}
		} else {
			if ui.IsTerminal() && !uninstallYes {
				fmt.Println("=== Uninstall PostgreSQL ===")
				fmt.Println()
				ans, err := ui.Prompt("This may remove PostgreSQL and its data. Continue? [yes/no]: ")
				if err != nil {
					return err
				}
				if strings.ToLower(strings.TrimSpace(ans)) != "yes" {
					fmt.Println("Cancelled.")
					return nil
				}
				fmt.Println()
			} else if !ui.IsTerminal() && !uninstallYes {
				return fmt.Errorf("Error: uninstall in non-interactive mode requires --yes")
			}
		}

		osName := sys.DetectOS()
		fmt.Println("Aniqlangan OS:", osName)

		var results []stepResult
		switch osName {
		case "mac":
			results = uninstallMac(uninstallHard)
		case "debian":
			results = uninstallDebian(uninstallHard)
		default:
			return fmt.Errorf("bu OS uchun avtomatik uninstall yo'q")
		}

		fmt.Println("\n=== Uninstall Summary ===")
		for _, r := range results {
			if r.label == "" {
				continue
			}
			status := "OK"
			if r.err != nil {
				status = "failed (" + r.err.Error() + ")"
			}
			fmt.Printf("- %s: %s\n", r.label, status)
		}

		fmt.Println("\nIf issues remain, you may need to manually clean lingering files or services.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolVar(&uninstallHard, "hard", false, "Perform a hard uninstall (stop services, uninstall packages, remove data/logs)")
	uninstallCmd.Flags().BoolVar(&uninstallYes, "yes", false, "Automatic yes to prompts (use with caution)")
}
