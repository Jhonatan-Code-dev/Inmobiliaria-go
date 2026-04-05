package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestMySQLUTCIntegration(t *testing.T) {
	t.Parallel()

	baseDSN := os.Getenv("TEST_MYSQL_DSN")
	if baseDSN == "" {
		t.Skip("TEST_MYSQL_DSN no definido; se omite prueba de integracion MySQL")
	}

	testDSN, dbName, cleanup := prepareTestMySQLDSN(t, baseDSN)
	defer cleanup()

	client, err := Setup(testDSN)
	if err != nil {
		t.Fatalf("Setup() error = %v", err)
	}
	defer client.Close()

	marca := time.Date(2026, 4, 3, 22, 30, 15, 0, time.UTC)
	if _, err := client.Empresa.
		Create().
		SetNombre("Empresa UTC Test").
		SetMoneda("USD").
		SetCreadoEn(marca).
		Save(t.Context()); err != nil {
		t.Fatalf("create empresa error = %v", err)
	}

	rawDSN := normalizeMySQLDSN(testDSN)
	db, err := sql.Open("mysql", rawDSN)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer db.Close()

	var creadoEn time.Time
	row := db.QueryRow("SELECT creado_en FROM empresas WHERE nombre = ?", "Empresa UTC Test")
	if err := row.Scan(&creadoEn); err != nil {
		t.Fatalf("scan error = %v", err)
	}

	if !creadoEn.Equal(marca) {
		t.Fatalf("creado_en = %s, want %s", creadoEn.Format(time.RFC3339), marca.Format(time.RFC3339))
	}
	if creadoEn.Location().String() != "UTC" {
		t.Fatalf("creado_en location = %s, want UTC", creadoEn.Location().String())
	}

	t.Logf("verified UTC read/write end-to-end on database %s", dbName)
}

func prepareTestMySQLDSN(t *testing.T, dsn string) (string, string, func()) {
	t.Helper()

	parts := strings.SplitN(dsn, "/", 2)
	if len(parts) != 2 {
		t.Fatalf("TEST_MYSQL_DSN invalido: %s", dsn)
	}

	dbSuffix := fmt.Sprintf("rentals_go_test_%d", time.Now().UTC().UnixNano())
	query := ""
	if idx := strings.Index(parts[1], "?"); idx >= 0 {
		query = parts[1][idx:]
	}
	testDSN := parts[0] + "/" + dbSuffix + query

	cleanup := func() {
		baseConn := parts[0] + "/"
		db, err := sql.Open("mysql", normalizeMySQLDSN(baseConn))
		if err != nil {
			t.Fatalf("cleanup sql.Open() error = %v", err)
		}
		defer db.Close()
		if _, err := db.Exec("DROP DATABASE IF EXISTS " + quoteIdentifier(dbSuffix)); err != nil {
			t.Fatalf("cleanup DROP DATABASE error = %v", err)
		}
	}

	return testDSN, dbSuffix, cleanup
}
