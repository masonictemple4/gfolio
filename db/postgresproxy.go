package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"cloud.google.com/go/cloudsqlconn"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgOptionsFn func(*gorm.DB)

func WithExistingConn() func(*gorm.DB) {
	conn, err := connectWithConnector()
	if err != nil {
		log.Fatalf("[with existing conn] %v\n", err)
	}
	return func(d *gorm.DB) {
		d, err = gorm.Open(postgres.New(postgres.Config{
			DriverName: "cloudsqlpostgres",
			Conn:       conn,
		}))

		if err != nil {
			log.Fatalf("[new postgres proxy] %v\n", err)
		}
	}
}

func WithDSN() func(*gorm.DB) {
	var (
		dbUser                 = promiseEnv("DB_USER")
		dbPwd                  = promiseEnv("DB_PASS")
		dbName                 = promiseEnv("DB_NAME")
		instanceConnectionName = promiseEnv("INSTANCE_CONNECTION_NAME")
	)

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", instanceConnectionName, dbUser, dbName, dbPwd)

	return func(d *gorm.DB) {
		d, err := gorm.Open(postgres.New(postgres.Config{
			DriverName: "cloudsqlpostgres",
			DSN:        dsn,
		}))

		if err != nil {
			log.Fatalf("[new postgres proxy] %v\n", err)
		}
	}
}

func NewPostgresGCPProxy(opts ...PgOptionsFn) *gorm.DB {
	var db *gorm.DB

	for _, opt := range opts {
		opt(db)
	}

	return db
}

func connectWithConnector() (*sql.DB, error) {
	var (
		dbUser                 = promiseEnv("DB_USER")
		dbPwd                  = promiseEnv("DB_PASS")
		dbName                 = promiseEnv("DB_NAME")
		instanceConnectionName = promiseEnv("INSTANCE_CONNECTION_NAME")
		usePrivate             = os.Getenv("PRIVATE_IP")
	)

	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)
	config, err := pgx.ParseConfig(dsn)

	if err != nil {
		return nil, err
	}

	var opts []cloudsqlconn.Option

	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}

	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)

	if err != nil {
		return nil, err
	}

	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName)
	}

	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open("pgx", dbURI)

	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	return dbPool, nil
}

func promiseEnv(name string) string {
	val, ok := os.LookupEnv(name)
	if !ok || val == "" {
		log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", name)
	}
	return val
}
