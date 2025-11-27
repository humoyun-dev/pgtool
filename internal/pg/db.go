package pg

import (
	"fmt"
	"strings"
)

func CreateDB(name, owner string) error {
	return runCreatedb(name, owner)
}

func DeleteDB(name string) error {
	return runDropdb(name)
}

func ListDBs() error {
	return runPsql("-c", `\l`)
}

func CreateUserAndDB(username, password, perms, dbName string) error {
	if err := CreateUser(username, password, perms); err != nil {
		return err
	}
	return CreateDB(dbName, username)
}

// ListDatabaseNames returns non-template database names for selection prompts.
func ListDatabaseNames() ([]string, error) {
	out, err := runPsqlOutput("-Atc", "SELECT datname FROM pg_database WHERE datistemplate = false ORDER BY datname;")
	if err != nil {
		if out != "" {
			return nil, fmt.Errorf("%w: %s", err, out)
		}
		return nil, err
	}

	lines := strings.Split(out, "\n")
	var dbs []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		dbs = append(dbs, line)
	}
	return dbs, nil
}
