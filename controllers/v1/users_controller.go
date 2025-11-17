package v1

import (
	"log"
	"net/http"

	"github.com/brandonrachal/gin-and-tonic/db"
	"github.com/brandonrachal/gin-and-tonic/models"
	"github.com/brandonrachal/gin-and-tonic/models/api"
	"github.com/gin-gonic/gin"
)

type UsersController struct {
	DBClient *db.Client
	logger   *log.Logger
}

func NewUsersController(logger *log.Logger, dbClient *db.Client) *UsersController {
	return &UsersController{
		DBClient: dbClient,
		logger:   logger,
	}
}

func (c *UsersController) CreateUserAction(ctx *gin.Context) {
	var user models.CreateUser
	if err := ctx.ShouldBindJSON(&user); err != nil {
		c.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, api.NewErrorMessage(err.Error()))
		return
	}
	result, resultErr := c.DBClient.CreateUser(ctx, user.FirstName, user.LastName, user.Email, user.Birthday.ToTime())
	if resultErr != nil {
		c.logger.Printf("Error inserting user - %s\n", resultErr)
		ctx.JSON(http.StatusInternalServerError, api.NewErrorMessage("Failed to insert user"))
		return
	}
	userId, userIdErr := result.LastInsertId()
	if userIdErr != nil {
		c.logger.Printf("Error getting the last id - %s\n", userIdErr)
		ctx.JSON(http.StatusInternalServerError, api.NewErrorMessage("Failed to insert user"))
		return
	}
	ctx.JSON(http.StatusOK, api.NewIdUserMessage(userId))
}

func (c *UsersController) GetUserAction(ctx *gin.Context) {
	var idUser models.IdUser
	if err := ctx.ShouldBindJSON(&idUser); err != nil {
		c.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, api.NewErrorMessage(err.Error()))
		return
	}
	user, userErr := c.DBClient.GetUser(ctx, idUser.Id)
	if userErr != nil {
		c.logger.Printf("Error retriving user id %d - %s\n", idUser.Id, userErr)
		ctx.JSON(http.StatusInternalServerError, api.NewErrorMessage("something went wrong"))
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (c *UsersController) UpdateUserAction(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		c.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, api.NewErrorMessage(err.Error()))
		return
	}
	_, resultErr := c.DBClient.UpdateUser(ctx, user.Id, user.FirstName, user.LastName, user.Email, user.Birthday.ToTime())
	if resultErr != nil {
		c.logger.Printf("Error inserting user - %s\n", resultErr)
		ctx.JSON(http.StatusInternalServerError, api.NewErrorMessage("Failed to insert user"))
		return
	}
	ctx.JSON(http.StatusOK, api.NewMessage("User updated successfully."))
}

func (c *UsersController) DeleteUserAction(ctx *gin.Context) {
	var user models.IdUser
	if err := ctx.ShouldBindJSON(&user); err != nil {
		c.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, api.NewErrorMessage(err.Error()))
		return
	}
	_, deleteUserErr := c.DBClient.DeleteUser(ctx, user.Id)
	if deleteUserErr != nil {
		c.logger.Printf("Error retrieving user id %d - %s\n", user.Id, deleteUserErr)
		ctx.JSON(http.StatusInternalServerError, api.NewErrorMessage("something went wrong"))
	}
	ctx.JSON(http.StatusOK, api.NewMessage("User deleted successfully."))
}

func (c *UsersController) GetUsersAction(ctx *gin.Context) {
	users, usersErr := c.DBClient.GetUsers(ctx)
	if usersErr != nil {
		c.logger.Printf("Error retrieving all users - %s\n", usersErr)
		ctx.JSON(http.StatusInternalServerError, api.NewErrorMessage("something went wrong"))
	}
	ctx.JSON(http.StatusOK, api.NewUsersMessage(users))
}

func (c *UsersController) GetUsersWithAgeAction(ctx *gin.Context) {
	usersWithAge, usersWithAgeErr := c.DBClient.GetUsersWithAge(ctx)
	if usersWithAgeErr != nil {
		c.logger.Printf("Error retrieving users with age - %s\n", usersWithAgeErr)
		ctx.JSON(http.StatusInternalServerError, api.NewErrorMessage("something went wrong"))
	}
	ctx.JSON(http.StatusOK, api.NewUsersWithAgeMessage(usersWithAge))
}

func (c *UsersController) GetAgeStatsAction(ctx *gin.Context) {
	ageStats, ageStatsErr := c.DBClient.GetAgeStats(ctx)
	if ageStatsErr != nil {
		c.logger.Printf("Error retrieving age stats - %s\n", ageStatsErr)
		ctx.JSON(http.StatusInternalServerError, api.NewErrorMessage("something went wrong"))
	}
	ctx.JSON(http.StatusOK, api.NewAgeStatsMessage(*ageStats))
}
