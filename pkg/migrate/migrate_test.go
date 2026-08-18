package migrate

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestParseName(t *testing.T) {
	tests := []struct {
		filename string
		version  int
		name     string
		isUp     bool
		ok       bool
	}{
		{"0001_create_users.up.sql", 1, "create_users", true, true},
		{"0002_add_index.down.sql", 2, "add_index", false, true},
		{"not_a_migration.txt", 0, "", false, false},
		{"0001_no_direction.sql", 0, "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			v, n, up, ok := parseName(tt.filename)
			if v != tt.version || n != tt.name || up != tt.isUp || ok != tt.ok {
				t.Fatalf("parseName(%q) = (%d,%q,%v,%v), want (%d,%q,%v,%v)",
					tt.filename, v, n, up, ok, tt.version, tt.name, tt.isUp, tt.ok)
			}
		})
	}
}

func TestPending(t *testing.T) {
	ms := []migration{{version: 1}, {version: 2}, {version: 3}}
	tests := []struct {
		name    string
		current int
		want    int
	}{
		{"全部待应用", 0, 3},
		{"部分待应用", 1, 2},
		{"已最新", 3, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(pending(ms, tt.current)); got != tt.want {
				t.Fatalf("pending() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestUpAppliesPending(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := New(db)
	r.ms = []migration{{version: 1, name: "init", up: "CREATE TABLE t (id INT)", down: "DROP TABLE t"}}

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT COALESCE").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(0))
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE t").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("INSERT INTO schema_migrations").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := r.Up(context.Background()); err != nil {
		t.Fatalf("Up() err = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("未满足的期望: %v", err)
	}
}

func TestUpIdempotent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := New(db)
	r.ms = []migration{{version: 1, name: "init", up: "CREATE TABLE t (id INT)"}}

	// 已应用版本 1：Up 只应查询版本，不产生任何事务。
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schema_migrations").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT COALESCE").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(1))

	if err := r.Up(context.Background()); err != nil {
		t.Fatalf("Up() err = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("幂等性未满足: %v", err)
	}
}
