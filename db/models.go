package db

import "github.com/brandonrachal/go-toolbox/jsonutils"

type UserId struct {
	Id int `db:"id" json:"id" form:"id" binding:"required"`
}

type NewUser struct {
	FirstName string               `db:"first_name" json:"first_name" form:"first_name" binding:"required"`
	LastName  string               `db:"last_name" json:"last_name" form:"last_name" binding:"required"`
	Email     string               `db:"email" json:"email" form:"email" binding:"required"`
	Birthday  jsonutils.SimpleDate `db:"birthday" json:"birthday" form:"birthday" binding:"required"`
}

type User struct {
	UserId
	NewUser
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
