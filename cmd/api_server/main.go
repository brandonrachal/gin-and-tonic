package main

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"time"

	"github.com/brandonrachal/gin-and-tonic/db"
	"github.com/brandonrachal/gin-and-tonic/server"
	"github.com/brandonrachal/go-toolbox/cliutils"
	"github.com/brandonrachal/go-toolbox/envutils"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	ctx, cancelFunc, logger := cliutils.InitDevConsole()
	defer cancelFunc()
	logger.SetOutput(gin.DefaultWriter)

	goEnv, goEnvErr := envutils.NewGoEnv()
	if goEnvErr != nil {
		logger.Fatalf("Could not retrieve the go env - %s\n", goEnvErr)
	}

	rootPath := goEnv.ModuleRootPath()
	sqliteDBPath := filepath.Join(rootPath, "data/sqlite_database.db")
	dbClient, dbClientErr := db.NewClient(sqliteDBPath)
	if dbClientErr != nil {
		logger.Fatalf("Could not retrieve the db client - %s\n", dbClientErr)
	}

	handler := server.NewHandler(logger, dbClient)
	router := gin.Default()
	router.GET("/ping", handler.Ping)
	router.POST("/user", handler.CreateUserHandler)
	router.GET("/user", handler.GetUserHandler)
	router.PUT("/user", handler.UpdateUserHandler)
	router.DELETE("/user", handler.DeleteUserHandler)
	router.GET("/users", handler.GetAllUsersHandler)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	logger.Println("Starting server on port 8080")
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("failed to start server - %s\n", err)
		}
	}()
	<-ctx.Done()
	cancelFunc()

	logger.Println("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown - %s\n", err)
	}
	logger.Println("Server exiting")
}
