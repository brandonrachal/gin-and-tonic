package internal

import (
	"fmt"
	"path/filepath"

	"github.com/brandonrachal/gin-and-tonic/db"
	"github.com/brandonrachal/go-toolbox/dbutils"
	"github.com/brandonrachal/go-toolbox/envutils"
	"github.com/brandonrachal/go-toolbox/migrations"
)

const (
	migrationTable = "goose_migrations"
	prodEnv        = "prod"
	devEnv         = "dev"
	testEnv        = "test"
)

var cachedGoEnv *envutils.GoEnv

func init() {
	var goEnvErr error
	cachedGoEnv, goEnvErr = envutils.NewGoEnv()
	if goEnvErr != nil {
		panic(fmt.Sprintf("Couldn't get the go env - %s\n", goEnvErr))
	}
}

func dbPath(env string) string {
	fileName := fmt.Sprintf("sqlite_%s_database.db", env)
	return filepath.Join(cachedGoEnv.ModuleRootPath(), "data", fileName)
}

func migrationsPath() string {
	return filepath.Join(cachedGoEnv.ModuleRootPath(), "migrations")
}

func dBMigrationClient(env string) (*migrations.Client, error) {
	return migrations.NewClient(dbutils.SQLite, dbPath(env), migrationTable, migrationsPath())
}

func ProdDBMigrationClient() (*migrations.Client, error) {
	return dBMigrationClient(prodEnv)
}

func ProdDBClient() (*db.Client, error) {
	return db.NewClient(dbPath(prodEnv))
}

func DevDBMigrationClient() (*migrations.Client, error) {
	return dBMigrationClient(devEnv)
}

func DevDBClient() (*db.Client, error) {
	return db.NewClient(dbPath(devEnv))
}

func TestDBMigrationClient() (*migrations.Client, error) {
	return dBMigrationClient(testEnv)
}

func TestDBClient() (*db.Client, error) {
	return db.NewClient(dbPath(testEnv))
}
