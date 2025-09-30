package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/brandonrachal/go-toolbox/cliutils"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancelFunc, logger := cliutils.InitDevConsole()
	defer cancelFunc()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome Gin Server",
		})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	cancelFunc()

	logger.Println("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}
	logger.Println("Server exiting")
}
