package pg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var errPsqlNotFound = errors.New("psql_not_found")

// defaultMetaDB returns the database to connect to for meta-queries (listing users/dbs).
func defaultMetaDB() string {
	if db := os.Getenv("PGTOOL_DB"); db != "" {
		return db
	}
	return "postgres"
}

// hasPostgresSystemUser reports whether a postgres system user exists (Linux only).
func hasPostgresSystemUser() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	cmd := exec.Command("id", "postgres")
	return cmd.Run() == nil
}

// psqlBin resolves the psql binary path with macOS Homebrew fallbacks.
func psqlBin() string {
	if p, err := exec.LookPath("psql"); err == nil {
		return p
	}
	if runtime.GOOS == "darwin" {
		candidates := []string{
			"/opt/homebrew/opt/postgresql@16/bin/psql",
			"/usr/local/opt/postgresql@16/bin/psql",
			"/opt/homebrew/opt/postgresql/bin/psql",
			"/usr/local/opt/postgresql/bin/psql",
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c
			}
		}
	}
	return "psql"
}

// createdbBin resolves the createdb binary.
func createdbBin() string {
	if p, err := exec.LookPath("createdb"); err == nil {
		return p
	}
	if runtime.GOOS == "darwin" {
		candidates := []string{
			"/opt/homebrew/opt/postgresql@16/bin/createdb",
			"/usr/local/opt/postgresql@16/bin/createdb",
			"/opt/homebrew/opt/postgresql/bin/createdb",
			"/usr/local/opt/postgresql/bin/createdb",
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c
			}
		}
	}
	return "createdb"
}

// dropdbBin resolves the dropdb binary.
func dropdbBin() string {
	if p, err := exec.LookPath("dropdb"); err == nil {
		return p
	}
	if runtime.GOOS == "darwin" {
		candidates := []string{
			"/opt/homebrew/opt/postgresql@16/bin/dropdb",
			"/usr/local/opt/postgresql@16/bin/dropdb",
			"/opt/homebrew/opt/postgresql/bin/dropdb",
			"/usr/local/opt/postgresql/bin/dropdb",
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c
			}
		}
	}
	return "dropdb"
}

// ensureDBArg prepends a -d defaultMetaDB() when not provided.
func ensureDBArg(args []string) ([]string, string) {
	hasDB := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-d" || arg == "--dbname" {
			hasDB = true
			break
		}
		if strings.HasPrefix(arg, "-d") && arg != "-d" {
			hasDB = true
			break
		}
		if strings.HasPrefix(arg, "--dbname=") {
			hasDB = true
			break
		}
	}

	if hasDB {
		return args, ""
	}
	db := defaultMetaDB()
	return append([]string{"-d", db}, args...), db
}

// buildCmd builds an *exec.Cmd to run a postgres tool honoring OS rules.
func buildCmd(bin string, args ...string) *exec.Cmd {
	if hasPostgresSystemUser() {
		all := append([]string{"-u", "postgres", bin}, args...)
		return exec.Command("sudo", all...)
	}
	return exec.Command(bin, args...)
}

func classifyAndWrap(db string, output string) error {
	orig := strings.TrimSpace(output)
	var short string
	var hints []string

	lower := strings.ToLower(orig)

	targetDB := db
	if targetDB == "" {
		needle := `database "`
		if idx := strings.Index(lower, needle); idx >= 0 {
			rest := lower[idx+len(needle):]
			if end := strings.Index(rest, `"`); end > 0 {
				targetDB = rest[:end]
			}
		}
	}

	if targetDB != "" && strings.Contains(lower, "database") && strings.Contains(lower, "does not exist") {
		short = fmt.Sprintf(`database "%s" does not exist`, targetDB)
		hints = []string{
			"List databases: psql -l",
			`Create the default meta database: psql -d template1 -c "CREATE DATABASE postgres;"`,
			"Or set PGTOOL_DB to an existing database, e.g.: export PGTOOL_DB=mydb",
		}
	} else if strings.Contains(lower, "could not connect to server") {
		short = "could not connect to PostgreSQL server"
		if runtime.GOOS == "darwin" {
			hints = []string{
				"Check if PostgreSQL is installed: brew list postgresql@16",
				"Start the service: brew services start postgresql@16",
				`Ensure PATH includes psql: echo 'export PATH="/opt/homebrew/opt/postgresql@16/bin:$PATH"' >> ~/.zshrc`,
			}
		} else {
			hints = []string{
				"Check service status: sudo systemctl status postgresql",
				"Start service: sudo systemctl start postgresql",
			}
		}
	} else if strings.Contains(lower, "psql: command not found") || strings.Contains(lower, "executable file not found") {
		short = "psql not found"
		if runtime.GOOS == "darwin" {
			hints = []string{
				"Install PostgreSQL: brew install postgresql@16",
				`Ensure PATH includes psql: echo 'export PATH="/opt/homebrew/opt/postgresql@16/bin:$PATH"' >> ~/.zshrc`,
			}
		} else {
			hints = []string{
				"Install PostgreSQL: sudo apt-get install postgresql postgresql-contrib",
			}
		}
	} else if strings.Contains(lower, "fatal:") {
		short = "psql command failed"
	} else {
		short = "psql command failed"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Error: %s\n", short)
	if len(hints) > 0 {
		b.WriteString("\nHints:\n")
		for _, h := range hints {
			b.WriteString("  - " + h + "\n")
		}
	}
	if orig != "" {
		b.WriteString("\nOriginal error:\n  ")
		b.WriteString(strings.ReplaceAll(orig, "\n", "\n  "))
		b.WriteString("\n")
	}
	return errors.New(b.String())
}

func runPsql(args ...string) error {
	args, db := ensureDBArg(args)
	cmd := buildCmd(psqlBin(), args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return classifyAndWrap(db, string(output))
	}
	if len(output) > 0 {
		fmt.Print(string(output))
	}
	return nil
}

func runCreatedb(db, owner string) error {
	cmd := buildCmd(createdbBin(), "-O", owner, db)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDropdb(db string) error {
	cmd := buildCmd(dropdbBin(), db)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runPsqlOutput executes psql and returns combined stdout/stderr for parsing.
func runPsqlOutput(args ...string) (string, error) {
	args, db := ensureDBArg(args)
	cmd := buildCmd(psqlBin(), args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", classifyAndWrap(db, string(output))
	}
	return string(output), nil
}
