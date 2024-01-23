package db

import (
	"fmt"
	"log"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresGCPProxy() *gorm.DB {

	var (
		dbUser                 = promiseEnv("DB_USER")
		dbPwd                  = promiseEnv("DB_PASS")
		dbName                 = promiseEnv("DB_NAME")
		instanceConnectionName = promiseEnv("INSTANCE_CONNECTION_NAME")
	)

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", instanceConnectionName, dbUser, dbName, dbPwd)

	println("connecting with DSN: ", dsn)

	var err error
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "cloudsqlpostgres",
		DSN:        dsn,
	}))

	if err != nil {
		log.Fatalf("[new postgres proxy] %v\n", err)
	}

	return db
}

func promiseEnv(name string) string {
	val, ok := os.LookupEnv(name)
	if !ok || val == "" {
		log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.\n", name)
	}
	return val
}
