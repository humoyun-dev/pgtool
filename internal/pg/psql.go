package pg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var errPsqlNotFound = errors.New("psql_not_found")

func hasPostgresSystemUser() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	cmd := exec.Command("id", "postgres")
	return cmd.Run() == nil
}

func binaryPath(name string) (string, error) {
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}

	if runtime.GOOS == "darwin" {
		// Homebrew defaults; keep a few fallbacks for Intel/ARM.
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
	}

	return "", fmt.Errorf("%s not found", name)
}

func psqlBin() (string, error) {
	p, err := binaryPath("psql")
	if err != nil {
		return "", errPsqlNotFound
	}
	return p, nil
}

func createdbBin() (string, error) {
	return binaryPath("createdb")
}

func dropdbBin() (string, error) {
	return binaryPath("dropdb")
}

// defaultMetaDB returns the database to connect to for meta-queries (listing users/dbs).
func defaultMetaDB() string {
	if db := os.Getenv("PGTOOL_DB"); db != "" {
		return db
	}
	return "postgres"
}

// buildCmd builds an *exec.Cmd to run a postgres tool honoring OS rules.
func buildCmd(bin string, args ...string) *exec.Cmd {
	if hasPostgresSystemUser() {
		all := append([]string{"-u", "postgres", bin}, args...)
		return exec.Command("sudo", all...)
	}
	return exec.Command(bin, args...)
}

func runPsql(args ...string) error {
	bin, err := psqlBin()
	if err != nil {
		if errors.Is(err, errPsqlNotFound) {
			return fmt.Errorf("Error: psql not found. Make sure PostgreSQL is installed and in PATH.")
		}
		return err
	}

	cmd := buildCmd(bin, args...)
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		fmt.Print(string(output))
	}
	if err != nil {
		return fmt.Errorf("Error: psql not found or cannot connect. Make sure PostgreSQL is installed, running, and accessible.")
	}
	return nil
}

func runCreatedb(db, owner string) error {
	bin, err := createdbBin()
	if err != nil {
		return err
	}

	cmd := buildCmd(bin, "-O", owner, db)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDropdb(db string) error {
	bin, err := dropdbBin()
	if err != nil {
		return err
	}

	cmd := buildCmd(bin, db)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runPsqlOutput executes psql and returns combined stdout/stderr for parsing.
func runPsqlOutput(args ...string) (string, error) {
	bin, err := psqlBin()
	if err != nil {
		if errors.Is(err, errPsqlNotFound) {
			return "", fmt.Errorf("Error: psql not found. Make sure PostgreSQL is installed and in PATH.")
		}
		return "", err
	}

	cmd := buildCmd(bin, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(output)), fmt.Errorf("Error: psql not found or cannot connect. Make sure PostgreSQL is installed, running, and accessible.")
	}
	return string(output), nil
}
