package controllers

import (
	"log"
	"net/http"

	"github.com/brandonrachal/gin-and-tonic/controllers/v1"
	"github.com/brandonrachal/gin-and-tonic/db"
	"github.com/gin-gonic/gin"
)

func GetRouter(logger *log.Logger, dbClient *db.Client) *gin.Engine {
	router := gin.Default()
	// All root routes
	router.GET("/ping", Ping)
	// All v1.0 routes
	v1Router := router.Group("/v1.0")
	// User Controller
	userController := v1.NewUsersController(logger, dbClient)
	v1Router.POST("/user", userController.CreateUserAction)
	v1Router.GET("/user", userController.GetUserAction)
	v1Router.PUT("/user", userController.UpdateUserAction)
	v1Router.DELETE("/user", userController.DeleteUserAction)
	v1Router.GET("/users", userController.GetUsersAction)
	v1Router.GET("/users_with_age", userController.GetUsersWithAgeAction)
	v1Router.GET("/age_stats", userController.GetAgeStatsAction)
	return router
}

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
