package pg

import (
	"os"
	"strings"
	"testing"
)

func TestBuildConnInfoDefaults(t *testing.T) {
	os.Unsetenv("PGHOST")
	os.Unsetenv("PGPORT")

	c := BuildConnInfo("user", "pass", "db")
	if c.Host != "localhost" {
		t.Errorf("expected localhost host, got %q", c.Host)
	}
	if c.Port != "5432" {
		t.Errorf("expected 5432 port, got %q", c.Port)
	}
}

func TestBuildConnInfoEnv(t *testing.T) {
	os.Setenv("PGHOST", "pg.example.com")
	os.Setenv("PGPORT", "6543")
	defer os.Unsetenv("PGHOST")
	defer os.Unsetenv("PGPORT")

	c := BuildConnInfo("user", "pass", "db")
	if c.Host != "pg.example.com" {
		t.Errorf("expected host from env, got %q", c.Host)
	}
	if c.Port != "6543" {
		t.Errorf("expected port from env, got %q", c.Port)
	}
}

func TestFastAPIDSNEncoding(t *testing.T) {
	c := ConnInfo{
		Username: "user",
		Password: "p@ss word",
		Database: "db",
		Host:     "localhost",
		Port:     "5432",
	}
	dsn := c.FastAPIDSN()
	if !strings.Contains(dsn, "postgresql+asyncpg://") {
		t.Fatalf("expected asyncpg scheme, got %s", dsn)
	}
	if strings.Contains(dsn, "p@ss word") {
		t.Fatalf("expected password to be URL-encoded, got %s", dsn)
	}
	if !strings.Contains(dsn, "localhost:5432/db") {
		t.Fatalf("expected host/port/db in dsn, got %s", dsn)
	}
}

func TestDjangoConfig(t *testing.T) {
	c := ConnInfo{
		Username: "user",
		Password: "pass",
		Database: "db",
		Host:     "localhost",
		Port:     "5432",
	}
	cfg := c.DjangoConfig()
	expected := []string{
		`"ENGINE": "django.db.backends.postgresql"`,
		`"NAME": "db"`,
		`"USER": "user"`,
		`"PASSWORD": "pass"`,
		`"HOST": "localhost"`,
		`"PORT": "5432"`,
	}
	for _, e := range expected {
		if !strings.Contains(cfg, e) {
			t.Fatalf("expected config to contain %q, got:\n%s", e, cfg)
		}
	}
}
