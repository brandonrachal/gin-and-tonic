package models

import (
	"fmt"

	"github.com/brandonrachal/go-toolbox/jsonutils"
)

// Database Models

type IdUser struct {
	Id int64 `db:"id" json:"id" form:"id" binding:"required"`
}

func GetIdUser(id int64) *IdUser {
	return &IdUser{Id: id}
}

type CreateUser struct {
	FirstName string               `db:"first_name" json:"first_name" form:"first_name" binding:"required"`
	LastName  string               `db:"last_name" json:"last_name" form:"last_name" binding:"required"`
	Email     string               `db:"email" json:"email" form:"email" binding:"required"`
	Birthday  jsonutils.SimpleDate `db:"birthday" json:"birthday" form:"birthday" binding:"required"`
}

func GetCreateUser(firstName, lastName, email, birthday string) (*CreateUser, error) {
	var birthdayDate jsonutils.SimpleDate
	err := birthdayDate.UnmarshalJSON([]byte(birthday))
	if err != nil {
		return nil, err
	}
	return &CreateUser{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Birthday:  birthdayDate,
	}, nil
}

type User struct {
	IdUser
	CreateUser
}

func (u *User) String() string {
	return fmt.Sprintf("User: %d, \"%s\", \"%s\", \"%s\", \"%s\",", u.Id, u.FirstName, u.LastName, u.Email, u.Birthday.String())
}

type UserWithAge struct {
	User
	AgeInYears int `db:"age_in_years" json:"age_in_years" form:"age_in_years" binding:"required"`
}

type AgeStats struct {
	Preteen   int `db:"preteen" json:"preteen"`
	Teen      int `db:"teens" json:"teens"`
	Twenties  int `db:"twenties" json:"twenties"`
	Thirties  int `db:"thirties" json:"thirties"`
	Forties   int `db:"forties" json:"forties"`
	Fifties   int `db:"fifties" json:"fifties"`
	Sixties   int `db:"sixties" json:"sixties"`
	Seventies int `db:"seventies" json:"seventies"`
	Eighties  int `db:"eighties" json:"eighties"`
	Nineties  int `db:"nineties" json:"nineties"`
	Centurion int `db:"centurion" json:"centurion"`
}
