package pg

import (
	"fmt"
	"net/url"
	"os"
)

type ConnInfo struct {
	Username string
	Password string
	Database string
	Host     string
	Port     string
}

// BuildConnInfo fills defaults from environment variables.
// Host: PGHOST or localhost
// Port: PGPORT or 5432
func BuildConnInfo(username, password, database string) ConnInfo {
	host := os.Getenv("PGHOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PGPORT")
	if port == "" {
		port = "5432"
	}
	return ConnInfo{
		Username: username,
		Password: password,
		Database: database,
		Host:     host,
		Port:     port,
	}
}

// FastAPIDSN returns a SQLAlchemy/asyncpg style DSN with URL-encoded credentials.
// Example: postgresql+asyncpg://user:pass@host:port/db
func (c ConnInfo) FastAPIDSN() string {
	user := url.UserPassword(c.Username, c.Password).String()
	return fmt.Sprintf("postgresql+asyncpg://%s@%s:%s/%s", user, c.Host, c.Port, c.Database)
}

// DjangoConfig returns a Django DATABASES snippet.
func (c ConnInfo) DjangoConfig() string {
	return fmt.Sprintf(`DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.postgresql",
        "NAME": "%s",
        "USER": "%s",
        "PASSWORD": "%s",
        "HOST": "%s",
        "PORT": "%s",
    }
}`, c.Database, c.Username, c.Password, c.Host, c.Port)
}
