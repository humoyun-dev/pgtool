package pg

import (
	"strings"
	"testing"
)

func assertLoginCount(t *testing.T, sql string, expected int) {
	t.Helper()
	if cnt := strings.Count(strings.ToUpper(sql), "LOGIN"); cnt != expected {
		t.Fatalf("expected %d LOGIN tokens, got %d in %q", expected, cnt, sql)
	}
}

func TestBuildCreateRoleSQL_None(t *testing.T) {
	sql := BuildCreateRoleSQL("test", "pass", "")
	assertLoginCount(t, sql, 1)
	if strings.Contains(sql, "CREATEDB") {
		t.Fatalf("did not expect CREATEDB, got %q", sql)
	}
}

func TestBuildCreateRoleSQL_LoginOnly(t *testing.T) {
	sql := BuildCreateRoleSQL("test", "pass", "LOGIN")
	assertLoginCount(t, sql, 1)
}

func TestBuildCreateRoleSQL_CreateDB(t *testing.T) {
	sql := BuildCreateRoleSQL("test", "pass", "CREATEDB")
	assertLoginCount(t, sql, 1)
	if !strings.Contains(sql, "CREATEDB") {
		t.Fatalf("expected CREATEDB, got %q", sql)
	}
}

func TestBuildCreateRoleSQL_LoginCreateDB(t *testing.T) {
	sql := BuildCreateRoleSQL("test", "pass", "LOGIN CREATEDB")
	assertLoginCount(t, sql, 1)
	if !strings.Contains(sql, "CREATEDB") {
		t.Fatalf("expected CREATEDB, got %q", sql)
	}
}

func TestBuildCreateRoleSQL_LoginCreateDBCreateRole(t *testing.T) {
	sql := BuildCreateRoleSQL("test", "pass", "LOGIN CREATEDB CREATEROLE")
	assertLoginCount(t, sql, 1)
	if !strings.Contains(sql, "CREATEDB") || !strings.Contains(sql, "CREATEROLE") {
		t.Fatalf("expected CREATEDB and CREATEROLE, got %q", sql)
	}
}

func TestBuildCreateRoleSQL_Custom(t *testing.T) {
	sql := BuildCreateRoleSQL("test", "pass", "LOGIN REPLICATION")
	assertLoginCount(t, sql, 1)
	if !strings.Contains(sql, "REPLICATION") {
		t.Fatalf("expected REPLICATION, got %q", sql)
	}
}
