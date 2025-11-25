package pg

import (
	"os"
	"os/exec"
)

func hasPostgresSystemUser() bool {
	cmd := exec.Command("id", "postgres")
	return cmd.Run() == nil
}

func runPsql(args ...string) error {
	var cmd *exec.Cmd

	if hasPostgresSystemUser() {
		all := append([]string{"-u", "postgres", "psql"}, args...)
		cmd = exec.Command("sudo", all...)
	} else {
		cmd = exec.Command("psql", args...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCreatedb(db, owner string) error {
	var cmd *exec.Cmd

	if hasPostgresSystemUser() {
		cmd = exec.Command("sudo", "-u", "postgres", "createdb", "-O", owner, db)
	} else {
		cmd = exec.Command("createdb", "-O", owner, db)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDropdb(db string) error {
	var cmd *exec.Cmd

	if hasPostgresSystemUser() {
		cmd = exec.Command("sudo", "-u", "postgres", "dropdb", db)
	} else {
		cmd = exec.Command("dropdb", db)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
