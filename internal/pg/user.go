package pg

import (
	"fmt"
	"strings"
)

func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// normalizePerms strips redundant LOGIN tokens since the base template already includes LOGIN.
func normalizePerms(perms string) string {
	fields := strings.Fields(perms)
	var out []string
	for _, f := range fields {
		if strings.EqualFold(f, "LOGIN") {
			continue
		}
		out = append(out, f)
	}
	return strings.Join(out, " ")
}

// BuildCreateRoleSQL builds a CREATE ROLE statement with LOGIN included by default.
func BuildCreateRoleSQL(username, password, perms string) string {
	u := escapeSingleQuotes(username)
	p := escapeSingleQuotes(password)
	cleanPerms := strings.TrimSpace(normalizePerms(perms))

	permClause := ""
	if cleanPerms != "" {
		permClause = " " + cleanPerms
	}

	return fmt.Sprintf(`CREATE ROLE "%s" LOGIN PASSWORD '%s'%s;`, u, p, permClause)
}

func buildCreateUserSQL(username, password, perms string) string {
	u := escapeSingleQuotes(username)
	p := escapeSingleQuotes(password)
	cleanPerms := strings.TrimSpace(normalizePerms(perms))

	permClause := ""
	if cleanPerms != "" {
		permClause = " " + cleanPerms
	}
	return fmt.Sprintf(`
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '%s') THEN
      CREATE ROLE "%s" LOGIN PASSWORD '%s'%s;
   ELSE
      RAISE NOTICE 'Role %s already exists';
   END IF;
END
$$;`, u, username, p, permClause, u)
}

func buildResetUserPasswordSQL(username, password string) string {
	p := escapeSingleQuotes(password)
	return fmt.Sprintf(`ALTER ROLE "%s" WITH PASSWORD '%s';`, username, p)
}

func buildDeleteUserSQL(username string) string {
	return fmt.Sprintf(`DROP ROLE IF EXISTS "%s";`, username)
}

// TestBuildCreateUserSQL testlar uchun helper.
func TestBuildCreateUserSQL(username, password, perms string) string {
	return buildCreateUserSQL(username, password, perms)
}

// TestBuildResetUserPasswordSQL testlar uchun helper.
func TestBuildResetUserPasswordSQL(username, password string) string {
	return buildResetUserPasswordSQL(username, password)
}

// TestBuildDeleteUserSQL testlar uchun helper.
func TestBuildDeleteUserSQL(username string) string {
	return buildDeleteUserSQL(username)
}

func CreateUser(username, password, perms string) error {
	sql := buildCreateUserSQL(username, password, perms)
	return runPsql("-v", "ON_ERROR_STOP=1", "-c", sql)
}

func ResetUserPassword(username, password string) error {
	sql := buildResetUserPasswordSQL(username, password)
	return runPsql("-v", "ON_ERROR_STOP=1", "-c", sql)
}

func DeleteUser(username string) error {
	sql := buildDeleteUserSQL(username)
	return runPsql("-v", "ON_ERROR_STOP=1", "-c", sql)
}

func ListUsers() error {
	return runPsql("-c", `\du`)
}

// ListRoleNames returns non-system role names for selection prompts.
func ListRoleNames() ([]string, error) {
	out, err := runPsqlOutput("-Atc", "SELECT rolname FROM pg_roles WHERE rolname NOT LIKE 'pg_%' ORDER BY rolname;")
	if err != nil {
		if out != "" {
			return nil, fmt.Errorf("%w: %s", err, out)
		}
		return nil, err
	}

	lines := strings.Split(out, "\n")
	var roles []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		roles = append(roles, line)
	}
	return roles, nil
}
