package server

import (
	"log"
	"net/http"

	"github.com/brandonrachal/gin-and-tonic/db"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	DBClient *db.Client
	logger   *log.Logger
}

func NewHandler(logger *log.Logger, dbClient *db.Client) *Handler {
	return &Handler{
		DBClient: dbClient,
		logger:   logger,
	}
}

func (h *Handler) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) CreateUserHandler(ctx *gin.Context) {
	var newUser db.NewUser
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		h.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, resultErr := h.DBClient.CreateUser(ctx, newUser.FirstName, newUser.LastName,
		newUser.Email, newUser.Birthday.ToTime())
	if resultErr != nil {
		h.logger.Printf("Error inserting user - %s\n", resultErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}
	userId, userIdErr := result.LastInsertId()
	if userIdErr != nil {
		h.logger.Printf("Error getting the last id - %s\n", userIdErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": map[string]any{"id": userId}})
}

func (h *Handler) GetUserHandler(ctx *gin.Context) {
	var userId db.UserId
	if err := ctx.ShouldBindJSON(&userId); err != nil {
		h.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, userErr := h.DBClient.User(ctx, userId.Id)
	if userErr != nil {
		h.logger.Printf("Error retriving user id %d - %s\n", userId.Id, userErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *Handler) UpdateUserHandler(ctx *gin.Context) {
	var user db.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, resultErr := h.DBClient.UpdateUser(ctx, user.Id, user.FirstName, user.LastName, user.Email, user.Birthday.ToTime())
	if resultErr != nil {
		h.logger.Printf("Error inserting user - %s\n", resultErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}
	updatedUser, updatedUserErr := h.DBClient.User(ctx, user.Id)
	if updatedUserErr != nil {
		h.logger.Printf("Error retrieving user id %d - %s\n", user.Id, updatedUserErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"user": updatedUser})
}

func (h *Handler) DeleteUserHandler(ctx *gin.Context) {
	var userId db.UserId
	if err := ctx.ShouldBindJSON(&userId); err != nil {
		h.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, deleteUserErr := h.DBClient.DeleteUser(ctx, userId.Id)
	if deleteUserErr != nil {
		h.logger.Printf("Error retrieving user id %d - %s\n", userId.Id, deleteUserErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully."})
}

func (h *Handler) GetUsersHandler(ctx *gin.Context) {
	users, usersErr := h.DBClient.Users(ctx)
	if usersErr != nil {
		h.logger.Printf("Error retrieving all users - %s\n", usersErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *Handler) GetUsersWithAgeHandler(ctx *gin.Context) {
	usersWithAge, usersWithAgeErr := h.DBClient.UsersWithAge(ctx)
	if usersWithAgeErr != nil {
		h.logger.Printf("Error retrieving users with age - %s\n", usersWithAgeErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"users": usersWithAge})
}

func (h *Handler) GetAgeStatsHandler(ctx *gin.Context) {
	ageStats, ageStatsErr := h.DBClient.AgeStats(ctx)
	if ageStatsErr != nil {
		h.logger.Printf("Error retrieving age stats - %s\n", ageStatsErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, gin.H{"age_stats": ageStats})
}
