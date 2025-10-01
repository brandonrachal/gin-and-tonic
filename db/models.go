package db

type User struct {
	Id        int    `db:"id" json:"id" form:"id"`
	FirstName string `db:"first_name" json:"first_name" binding:"required"`
	LastName  string `db:"last_name" json:"last_name" binding:"required"`
	Email     string `db:"email" json:"email" binding:"required"`
}
