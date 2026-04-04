	oracleDSN := fmt.Sprintf(
		`user="%s" password="%s" connectString="%s"`,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBTNS,
	)



    package db

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"rentals-go/ent"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/godror/godror" // Driver para Oracle
)

type Oracle struct {
	dsn    string
	client *ent.Client
	db     *sql.DB
	once   sync.Once
	err    error
}

func NewOracle(dsn string) *Oracle {
	return &Oracle{dsn: dsn}
}

func (o *Oracle) GetClient() (*ent.Client, error) {
	o.once.Do(func() {
		// En Oracle ATP con TLS, el DSN que viene del config es la cadena larga
		sqlDB, err := sql.Open("godror", o.dsn)
		if err != nil {
			o.err = err
			return
		}

		// Configuración profesional del Pool
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(20) // ATP suele tener límites de procesos, 20 es seguro
		sqlDB.SetConnMaxLifetime(time.Hour)

		if err := sqlDB.Ping(); err != nil {
			o.err = err
			return
		}

		o.db = sqlDB
		// Usamos el dialecto "oracle" para Ent
		driver := entsql.OpenDB("oracle", sqlDB)
		o.client = ent.NewClient(ent.Driver(driver))

		log.Println("✅ Oracle Autonomous DB conectado correctamente (TLS)")
	})

	return o.client, o.err
}

func (o *Oracle) Ping(ctx context.Context) error {
	if o.db == nil {
		return sql.ErrConnDone
	}
	return o.db.PingContext(ctx)
}

func (o *Oracle) Migrate() error {
	if o.client == nil {
		return sql.ErrConnDone
	}
	if err := o.client.Schema.Create(context.Background()); err != nil {
		return err
	}
	log.Println("✅ Migraciones de Oracle ejecutadas correctamente")
	return nil
}

func (o *Oracle) Close() error {
	if o.client != nil {
		return o.client.Close()
	}
	return nil
}
