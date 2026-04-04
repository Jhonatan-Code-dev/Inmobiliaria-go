package db

import (
	"strings"
	"testing"
)

func TestNormalizeMySQLDSN(t *testing.T) {
	t.Parallel()

	dsn := "root:123456@tcp(localhost:3306)/alquileres"
	got := normalizeMySQLDSN(dsn)

	for _, expected := range []string{
		"parseTime=True",
		"loc=UTC",
		"time_zone=%27%2B00%3A00%27",
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("normalizeMySQLDSN() missing %s in %s", expected, got)
		}
	}
}

func TestNormalizePostgresDSN(t *testing.T) {
	t.Parallel()

	dsn := "postgres://user:pass@localhost:5432/app?sslmode=disable"
	got := normalizePostgresDSN(dsn)

	if !strings.Contains(got, "TimeZone=UTC") {
		t.Fatalf("normalizePostgresDSN() missing TimeZone=UTC in %s", got)
	}
}
