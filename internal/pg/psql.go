package pg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func hasPostgresSystemUser() bool {
	cmd := exec.Command("id", "postgres")
	return cmd.Run() == nil
}

// findLocalBinary attempts to resolve a binary path without sudo wrapping.
func findLocalBinary(name string) (string, error) {
	// First, look in PATH
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}

	switch runtime.GOOS {
	case "darwin":
		// Common Homebrew locations (especially Apple Silicon)
		candidates := []string{
			filepath.Join("/opt/homebrew/opt/postgresql@16/bin", name),
			filepath.Join("/opt/homebrew/opt/postgresql/bin", name),
			filepath.Join("/usr/local/opt/postgresql@16/bin", name),
			filepath.Join("/usr/local/opt/postgresql/bin", name),
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c, nil
			}
		}
	case "linux":
		candidates := []string{
			filepath.Join("/usr/bin", name),
			filepath.Join("/usr/local/bin", name),
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c, nil
			}
		}
	}
	return "", fmt.Errorf("binary %s not found", name)
}

var errPsqlNotFound = errors.New("psql_not_found")

// buildCmd builds an *exec.Cmd to run a postgres tool honoring OS rules.
func buildCmd(tool string, args ...string) (*exec.Cmd, error) {
	// Linux with postgres system user: use sudo wrapper
	if runtime.GOOS == "linux" && hasPostgresSystemUser() {
		all := append([]string{"-u", "postgres", tool}, args...)
		return exec.Command("sudo", all...), nil
	}

	// Otherwise, find local binary
	path, err := findLocalBinary(tool)
	if err != nil {
		if tool == "psql" {
			return nil, errPsqlNotFound
		}
		return nil, err
	}
	return exec.Command(path, args...), nil
}

func runPsql(args ...string) error {
	cmd, err := buildCmd("psql", args...)
	if err != nil {
		if errors.Is(err, errPsqlNotFound) {
			if runtime.GOOS == "darwin" {
				return fmt.Errorf("Error: psql topilmadi. PostgreSQL o'rnatilgan va PATH'ga qo'shilganiga ishonch hosil qil. macOS hint: /opt/homebrew/opt/postgresql@16/bin/psql")
			}
			return fmt.Errorf("Error: psql topilmadi. PostgreSQL o'rnatilgan va PATH'ga qo'shilganiga ishonch hosil qil.")
		}
		return err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCreatedb(db, owner string) error {
	cmd, err := buildCmd("createdb", "-O", owner, db)
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDropdb(db string) error {
	cmd, err := buildCmd("dropdb", db)
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
