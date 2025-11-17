package main

import (
	"context"
	"fmt"
	"os"

	"github.com/brandonrachal/gin-and-tonic/internal"
	"github.com/brandonrachal/go-toolbox/cliutils"
	"github.com/brandonrachal/go-toolbox/migrations"
)

func main() {
	ctx, cancelFunc := cliutils.InitSignals(context.Background())
	defer cancelFunc()

	migrateCmdData, migrateCmdDataErr := migrations.GetMigrationCmdData(false)
	if migrateCmdDataErr != nil {
		fmt.Printf("error getting migration cmd data - %s\n", migrateCmdDataErr)
		os.Exit(1)
	} else if migrateCmdData.ConfirmExit {
		fmt.Println("Confirmation failed. Skipping operation and exiting.")
		os.Exit(1)
	}

	migrateClient, migrateClientErr := internal.ProdDBMigrationClient()
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
