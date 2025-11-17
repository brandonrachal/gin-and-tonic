package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/brandonrachal/go-toolbox/jsonutils"
	"github.com/stretchr/testify/require"
)

const (
	IdUserJson     = `{"id": 1}`
	CreateUserJson = `{"first_name": "Brandon", "last_name": "Rachal", "email": "brandon.rachal@gmail.com", "birthday": "2025-10-12"}`
	UserJson       = `{"id": 1, "first_name": "Brandon", "last_name": "Rachal", "email": "brandon.rachal@gmail.com", "birthday": "2025-10-12"}`
)

func TestIdUser(t *testing.T) {
	req := require.New(t)
	var user IdUser
	err := json.Unmarshal([]byte(IdUserJson), &user)
	req.NoError(err)
	req.Equal(int64(1), user.Id)
}

func TestCreateUser(t *testing.T) {
	req := require.New(t)
	birthday, birthdayErr := time.Parse(jsonutils.SimpleDateFormat, "2025-10-12")
	req.NoError(birthdayErr)
	var user CreateUser
	err := json.Unmarshal([]byte(CreateUserJson), &user)
	req.NoError(err)
	req.Equal("Brandon", user.FirstName)
	req.Equal("Rachal", user.LastName)
	req.Equal("brandon.rachal@gmail.com", user.Email)
	req.Equal(birthday, user.Birthday.ToTime())
}

func TestUser(t *testing.T) {
	req := require.New(t)
	birthday, birthdayErr := time.Parse(jsonutils.SimpleDateFormat, "2025-10-12")
	req.NoError(birthdayErr)
	var user User
	err := json.Unmarshal([]byte(UserJson), &user)
	req.NoError(err)
	req.Equal(int64(1), user.Id)
	req.Equal("Brandon", user.FirstName)
	req.Equal("Rachal", user.LastName)
	req.Equal("brandon.rachal@gmail.com", user.Email)
	req.Equal(birthday, user.Birthday.ToTime())
}
