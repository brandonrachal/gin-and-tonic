package db

type UserId struct {
	Id int `db:"id" json:"id" form:"id" binding:"required"`
}

type NewUser struct {
	FirstName string `db:"first_name" json:"first_name" form:"first_name" binding:"required"`
	LastName  string `db:"last_name" json:"last_name" form:"last_name" binding:"required"`
	Email     string `db:"email" json:"email" form:"email" binding:"required"`
}

type User struct {
	UserId
	NewUser
}
