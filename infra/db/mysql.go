package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"rentals-go/ent"
	"rentals-go/ent/migrate"

	entsql "entgo.io/ent/dialect/sql"

	// Drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type DB struct {
	dsn    string
	client *ent.Client
	db     *sql.DB
	once   sync.Once
	err    error
	driver string
}

// NewDB crea un cliente que detecta el driver por la URL
func NewDB(dsn string) *DB {
	driver := detectDriver(dsn)
	dsn = normalizeDSN(dsn, driver)
	return &DB{
		dsn:    dsn,
		driver: driver,
	}
}

// detectDriver decide el driver según la DSN
func detectDriver(dsn string) string {
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		return "postgres"
	}
	// Por defecto mysql
	return "mysql"
}

func normalizeDSN(dsn, driver string) string {
	switch driver {
	case "postgres":
		return normalizePostgresDSN(dsn)
	default:
		return normalizeMySQLDSN(dsn)
	}
}

func normalizePostgresDSN(dsn string) string {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return dsn
	}
	query := parsed.Query()
	if query.Get("TimeZone") == "" && query.Get("timezone") == "" {
		query.Set("TimeZone", "UTC")
	}
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func normalizeMySQLDSN(dsn string) string {
	base := dsn
	rawQuery := ""
	if idx := strings.Index(dsn, "?"); idx >= 0 {
		base = dsn[:idx]
		rawQuery = dsn[idx+1:]
	}

	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return dsn
	}
	if values.Get("parseTime") == "" {
		values.Set("parseTime", "True")
	}
	if values.Get("loc") == "" {
		values.Set("loc", "UTC")
	}
	if values.Get("time_zone") == "" {
		// Fuerza la sesion MySQL/MariaDB a UTC sin depender de la configuracion del servidor.
		values.Set("time_zone", "'+00:00'")
	}

	encoded := values.Encode()
	if encoded == "" {
		return base
	}
	return base + "?" + encoded
}

func (dbClient *DB) ensureDatabaseExists() error {
	// Parseamos la DSN para MySQL (formato: user:pass@tcp(host:port)/dbname)
	parts := strings.Split(dbClient.dsn, "/")
	if len(parts) < 2 {
		return nil // No se pudo parsear el nombre de la BD
	}

	baseDSN := parts[0] + "/"
	dbName := strings.Split(parts[1], "?")[0]

	// Conectamos sin base de datos
	db, err := sql.Open("mysql", baseDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", quoteIdentifier(dbName)))
	return err
}

// quoteIdentifier escapa el nombre para prevenir caracteres inválidos en SQL.
func quoteIdentifier(name string) string {
	if name == "" {
		return ""
	}
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

func (dbClient *DB) GetClient() (*ent.Client, error) {
	dbClient.once.Do(func() {
		if dbClient.driver == "mysql" {
			if err := dbClient.ensureDatabaseExists(); err != nil {
				dbClient.err = err
				return
			}
		}

		sqlDB, err := sql.Open(dbClient.driver, dbClient.dsn)
		if err != nil {
			dbClient.err = err
			return
		}

		// Pool
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		// Probar conexión
		if err := sqlDB.Ping(); err != nil {
			dbClient.err = err
			return
		}

		dbClient.db = sqlDB
		driver := entsql.OpenDB(dbClient.driver, sqlDB)
		dbClient.client = ent.NewClient(ent.Driver(driver))

		log.Printf("🟦 %s conectado correctamente\n", strings.Title(dbClient.driver))
	})

	return dbClient.client, dbClient.err
}

func (dbClient *DB) Ping(ctx context.Context) error {
	if dbClient.db == nil {
		return sql.ErrConnDone
	}
	return dbClient.db.PingContext(ctx)
}

func (dbClient *DB) Migrate() error {
	if dbClient.client == nil {
		return sql.ErrConnDone
	}

	if err := dbClient.client.Schema.Create(
		context.Background(),
		migrate.WithDropColumn(true),
		migrate.WithDropIndex(true),
	); err != nil {
		return err
	}

	log.Printf("🟦 Migraciones ejecutadas correctamente en %s\n", strings.Title(dbClient.driver))
	return nil
}

func (dbClient *DB) Close() error {
	if dbClient.client != nil {
		return dbClient.client.Close()
	}
	return nil
}
