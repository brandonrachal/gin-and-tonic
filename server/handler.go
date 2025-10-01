package server

import (
	"log"
	"net/http"

	"github.com/brandonrachal/gin-and-tonic/db"
	"github.com/brandonrachal/go-toolbox/datautils"
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

func (h *Handler) GetUserHandler(ctx *gin.Context) {
	formUserId, exists := ctx.GetPostForm("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
	}
	userId, userIdErr := datautils.ToInt(formUserId)
	if userIdErr != nil {
		h.logger.Printf("Error converting user id to int - %s\n", userIdErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	user, userErr := h.DBClient.GetUser(ctx, *userId)
	if userErr != nil {
		h.logger.Printf("Error retriving user id '%s' - %s\n", formUserId, userIdErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) InsertUserHandler(ctx *gin.Context) {
	var user db.User
	if err := ctx.ShouldBind(&user); err != nil {
		h.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, resultErr := h.DBClient.InsertUser(ctx, user.FirstName, user.LastName, user.Email)
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
	resp := gin.H{"user_id": userId}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) UpdateUserHandler(ctx *gin.Context) {
	var user db.User
	if err := ctx.ShouldBind(&user); err != nil {
		h.logger.Printf("Error binding user - %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, resultErr := h.DBClient.UpdateUser(ctx, user.Id, user.FirstName, user.LastName, user.Email)
	if resultErr != nil {
		h.logger.Printf("Error inserting user - %s\n", resultErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}
	resp := gin.H{"user_id": user.Id}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteUserHandler(ctx *gin.Context) {
	formUserId, exists := ctx.GetPostForm("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
	}
	userId, userIdErr := datautils.ToInt(formUserId)
	if userIdErr != nil {
		h.logger.Printf("Error converting user id to int - %s\n", userIdErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	user, userErr := h.DBClient.GetUser(ctx, *userId)
	if userErr != nil {
		h.logger.Printf("Error retriving user id '%s' - %s\n", formUserId, userIdErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) GetAllUsersHandler(ctx *gin.Context) {
	users, usersErr := h.DBClient.GetAllUsers(ctx)
	if usersErr != nil {
		h.logger.Printf("Error retriving all users - %s\n", usersErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
	ctx.JSON(http.StatusOK, users)
}
