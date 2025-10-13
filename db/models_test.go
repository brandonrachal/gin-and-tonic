package db

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/brandonrachal/go-toolbox/jsonutils"
	"github.com/stretchr/testify/require"
)

const (
	UserIdJson  = `{"id": 1}`
	NewUserJson = `{"first_name": "Brandon", "last_name": "Rachal", "email": "brandon.rachal@gmail.com", "birthday": "2025-10-12"}`
	UserJson    = `{"id": 1, "first_name": "Brandon", "last_name": "Rachal", "email": "brandon.rachal@gmail.com", "birthday": "2025-10-12"}`
)

func TestUserId(t *testing.T) {
	req := require.New(t)
	var userId UserId
	err := json.Unmarshal([]byte(UserIdJson), &userId)
	req.NoError(err)
	req.Equal(1, userId.Id)
}

func TestNewUser(t *testing.T) {
	req := require.New(t)
	birthday, birthdayErr := time.Parse(jsonutils.SimpleDateFormat, "2025-10-12")
	req.NoError(birthdayErr)
	var newUser NewUser
	err := json.Unmarshal([]byte(NewUserJson), &newUser)
	req.NoError(err)
	req.Equal("Brandon", newUser.FirstName)
	req.Equal("Rachal", newUser.LastName)
	req.Equal("brandon.rachal@gmail.com", newUser.Email)
	req.Equal(birthday, newUser.Birthday.ToTime())
}

func TestUser(t *testing.T) {
	req := require.New(t)
	birthday, birthdayErr := time.Parse(jsonutils.SimpleDateFormat, "2025-10-12")
	req.NoError(birthdayErr)
	var user User
	err := json.Unmarshal([]byte(UserJson), &user)
	req.NoError(err)
	req.Equal(1, user.Id)
	req.Equal("Brandon", user.FirstName)
	req.Equal("Rachal", user.LastName)
	req.Equal("brandon.rachal@gmail.com", user.Email)
	req.Equal(birthday, user.Birthday.ToTime())
}
