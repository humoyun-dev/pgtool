package pg

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
