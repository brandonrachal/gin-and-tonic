package main

import (
	"context"
	"fmt"
	"os"

	"github.com/brandonrachal/go-toolbox/cliutils"
	"github.com/brandonrachal/go-toolbox/dbutils"
	"github.com/brandonrachal/go-toolbox/envutils"
	"github.com/brandonrachal/go-toolbox/migrations"
)

const (
	migrationTable = "goose_migrations"
)

func main() {
	bgCtx := context.Background()
	ctx, cancelFunc := cliutils.InitSignals(bgCtx)
	defer cancelFunc()

	goEnv, goEnvErr := envutils.NewGoEnv()
	if goEnvErr != nil {
		fmt.Printf("error loading go enviroment - %s", goEnvErr)
		os.Exit(1)
	}
	moduleRoot := goEnv.ModuleRootPath()
	migrationDir := fmt.Sprintf("%s/migrations", moduleRoot)
	dsn := fmt.Sprintf("%s/data/sqlite_database.db", moduleRoot)

	migrateCmdData, migrateCmdDataErr := migrations.GetMigrationCmdData(false)
	if migrateCmdDataErr != nil {
		fmt.Printf("error getting migration cmd data - %s\n", migrateCmdDataErr)
		os.Exit(1)
	} else if migrateCmdData.ConfirmExit {
		fmt.Println("Confirmation failed. Skipping operation and exiting.")
		os.Exit(1)
	}

	migrateClient, migrateClientErr := migrations.NewClient(dbutils.SQLite, dsn, migrationTable, migrationDir)
	if migrateClientErr != nil {
		fmt.Printf("error getting new migration client - %s\n", migrateClientErr)
		os.Exit(1)
	}

	runMigrationCmdErr := migrations.RunMigrationCmd(ctx, migrateCmdData, migrateClient)
	if runMigrationCmdErr != nil {
		fmt.Printf("error running migration command - %s\n", runMigrationCmdErr)
		os.Exit(1)
	}

	fmt.Println("Command completed successfully")
}
