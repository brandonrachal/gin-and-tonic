package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brandonrachal/gin-and-tonic/controllers"
	"github.com/brandonrachal/gin-and-tonic/internal"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	logger.SetOutput(gin.DefaultWriter)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancelFunc()
	}()

	dbClient, dbClientErr := internal.ProdDBClient()
	if dbClientErr != nil {
		logger.Fatalf("Could not retrieve the db client - %s\n", dbClientErr)
	}

	router := controllers.GetRouter(logger, dbClient)
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
