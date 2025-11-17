package api

import "github.com/brandonrachal/gin-and-tonic/models"

type Message struct {
	Message string `json:"message"`
}

func NewMessage(message string) *Message {
	return &Message{
		Message: message,
	}
}

type ErrorMessage struct {
	Error string `json:"error"`
}

func NewErrorMessage(error string) ErrorMessage {
	return ErrorMessage{Error: error}
}

type IdUserMessage struct {
	User models.IdUser `json:"user"`
}

func NewIdUserMessage(id int64) IdUserMessage {
	modelUserId := models.GetIdUser(id)
	return IdUserMessage{
		User: *modelUserId,
	}
}

type UsersMessage struct {
	Users []models.User `json:"users"`
}

func NewUsersMessage(users []models.User) UsersMessage {
	return UsersMessage{
		Users: users,
	}
}

type UsersWithAgeMessage struct {
	Users []models.UserWithAge `json:"users"`
}

func NewUsersWithAgeMessage(users []models.UserWithAge) UsersWithAgeMessage {
	return UsersWithAgeMessage{
		Users: users,
	}
}

type AgeStatsMessage struct {
	AgeStats models.AgeStats `json:"age_stats"`
}

func NewAgeStatsMessage(ageStats models.AgeStats) AgeStatsMessage {
	return AgeStatsMessage{
		AgeStats: ageStats,
	}
}
