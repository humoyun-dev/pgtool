package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/humoyun-dev/pgtool/internal/pg"
	"github.com/humoyun-dev/pgtool/internal/ui"
	"github.com/spf13/cobra"
)

var (
	createUserName string
	createUserPass string
	createUserPerm string
)

const defaultPerms = "LOGIN CREATEDB"

func choosePermissions(current string, interactive bool) (string, error) {
	if current != "" {
		return current, nil
	}

	if !interactive {
		// Non-interactive defaults to a safe baseline.
		return defaultPerms, nil
	}

	options := []struct {
		label string
		value string
	}{
		{"NONE", ""},
		{"LOGIN", "LOGIN"},
		{"CREATEDB", "CREATEDB"},
		{"CREATEROLE", "CREATEROLE"},
		{"REPLICATION", "REPLICATION"},
		{"SUPERUSER", "SUPERUSER"},
		{"LOGIN + CREATEDB", "LOGIN CREATEDB"},
		{"LOGIN + CREATEDB + CREATEROLE", "LOGIN CREATEDB CREATEROLE"},
		{"Custom (enter manually)", ""},
	}

	fmt.Println("=== Permissions ===")
	for i, opt := range options {
		fmt.Printf("  %d) %s\n", i+1, opt.label)
	}
	fmt.Print("Select an option [number, empty = 7]: ")

	choiceRaw, err := ui.Prompt("")
	if err != nil {
		return "", err
	}
	choiceRaw = strings.TrimSpace(choiceRaw)
	if choiceRaw == "" {
		choiceRaw = "7"
	}

	choice, err := strconv.Atoi(choiceRaw)
	if err != nil || choice < 1 || choice > len(options) {
		return "", fmt.Errorf("Error: invalid permissions selection.")
	}

	if choice == len(options) {
		custom, err := ui.Prompt("Enter custom permissions (raw SQL fragment, e.g. 'LOGIN CREATEDB'): ")
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(custom), nil
	}

	return options[choice-1].value, nil
}

var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "PostgreSQL user yaratish",
	RunE: func(cmd *cobra.Command, args []string) error {
		interactive := ui.IsTerminal()

		if interactive {
			fmt.Println("=== Create User ===")
			fmt.Println()
		}

		if interactive {
			if createUserName == "" {
				v, err := ui.Prompt("Username: ")
				if err != nil {
					return err
				}
				createUserName = v
			}
			if createUserPass == "" {
				v, err := ui.Prompt("Password: ")
				if err != nil {
					return err
				}
				createUserPass = v
			}
		} else if createUserName == "" || createUserPass == "" {
			return fmt.Errorf("Error: username and password are required (flags or interactive input).")
		}

		perms, err := choosePermissions(createUserPerm, interactive)
		if err != nil {
			return err
		}
		createUserPerm = perms

		if createUserName == "" || createUserPass == "" {
			return fmt.Errorf("Error: username and password are required (flags or interactive input).")
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
