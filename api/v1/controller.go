package v1

import (
	"log"
	"net/http"

	"github.com/brandonrachal/gin-and-tonic/db"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	DBClient *db.Client
	logger   *log.Logger
}

func NewController(logger *log.Logger, dbClient *db.Client) *Controller {
	return &Controller{
		DBClient: dbClient,
		logger:   logger,
	}
}

func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (c *Controller) CreateUserAction(ctx *gin.Context) {
	var newUser db.NewUser
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		c.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, resultErr := c.DBClient.CreateUser(ctx, newUser.FirstName, newUser.LastName,
		newUser.Email, newUser.Birthday.ToTime())
	if resultErr != nil {
		c.logger.Printf("Error inserting user - %s\n", resultErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}
	userId, userIdErr := result.LastInsertId()
	if userIdErr != nil {
		c.logger.Printf("Error getting the last id - %s\n", userIdErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": map[string]any{"id": userId}})
}

func (c *Controller) GetUserAction(ctx *gin.Context) {
	var userId db.UserId
	if err := ctx.ShouldBindJSON(&userId); err != nil {
		c.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, userErr := c.DBClient.User(ctx, userId.Id)
	if userErr != nil {
		c.logger.Printf("Error retriving user id %d - %s\n", userId.Id, userErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (c *Controller) UpdateUserAction(ctx *gin.Context) {
	var user db.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		c.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, resultErr := c.DBClient.UpdateUser(ctx, user.Id, user.FirstName, user.LastName, user.Email, user.Birthday.ToTime())
	if resultErr != nil {
		c.logger.Printf("Error inserting user - %s\n", resultErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}
	updatedUser, updatedUserErr := c.DBClient.User(ctx, user.Id)
	if updatedUserErr != nil {
		c.logger.Printf("Error retrieving user id %d - %s\n", user.Id, updatedUserErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"user": updatedUser})
}

func (c *Controller) DeleteUserAction(ctx *gin.Context) {
	var userId db.UserId
	if err := ctx.ShouldBindJSON(&userId); err != nil {
		c.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, deleteUserErr := c.DBClient.DeleteUser(ctx, userId.Id)
	if deleteUserErr != nil {
		c.logger.Printf("Error retrieving user id %d - %s\n", userId.Id, deleteUserErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully."})
}

func (c *Controller) GetUsersAction(ctx *gin.Context) {
	users, usersErr := c.DBClient.Users(ctx)
	if usersErr != nil {
		c.logger.Printf("Error retrieving all users - %s\n", usersErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (c *Controller) GetUsersWithAgeAction(ctx *gin.Context) {
	usersWithAge, usersWithAgeErr := c.DBClient.UsersWithAge(ctx)
	if usersWithAgeErr != nil {
		c.logger.Printf("Error retrieving users with age - %s\n", usersWithAgeErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"users": usersWithAge})
}

func (c *Controller) GetAgeStatsAction(ctx *gin.Context) {
	ageStats, ageStatsErr := c.DBClient.AgeStats(ctx)
	if ageStatsErr != nil {
		c.logger.Printf("Error retrieving age stats - %s\n", ageStatsErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"age_stats": ageStats})
}
