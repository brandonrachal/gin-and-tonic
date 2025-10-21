package v1

import (
	"log"

	"github.com/brandonrachal/gin-and-tonic/db"
	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine, logger *log.Logger, dbClient *db.Client) {
	controller := NewController(logger, dbClient)
	v1Router := router.Group("/v1.0")
	v1Router.GET("/ping", controller.Ping)
	v1Router.POST("/user", controller.CreateUserAction)
	v1Router.GET("/user", controller.GetUserAction)
	v1Router.PUT("/user", controller.UpdateUserAction)
	v1Router.DELETE("/user", controller.DeleteUserAction)
	v1Router.GET("/users", controller.GetUsersAction)
	v1Router.GET("/users_with_age", controller.GetUsersWithAgeAction)
	v1Router.GET("/age_stats", controller.GetAgeStatsAction)
}
