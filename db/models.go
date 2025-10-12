package db

import (
	"strings"
	"time"
)

const DateFormat = "2006-01-02"

type Date time.Time

func (d *Date) UnmarshalJSON(jsonBytes []byte) error {
	s := strings.Trim(string(jsonBytes), `"`) // Remove quotes
	if s == "null" || s == "" {
		*d = Date(time.Time{}) // Handle null or empty string
		return nil
	}
	date, dateErr := time.Parse(DateFormat, s)
	if dateErr != nil {
		return dateErr
	}
	*d = Date(date)
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	t := time.Time(*d)
	formatted := t.Format(DateFormat)
	jsonStr := "\"" + formatted + "\""
	return []byte(jsonStr), nil
}

func (d *Date) ToTime() time.Time {
	return time.Time(*d)
}

type UserId struct {
	Id int `db:"id" json:"id" form:"id" binding:"required"`
}

type NewUser struct {
	FirstName string `db:"first_name" json:"first_name" form:"first_name" binding:"required"`
	LastName  string `db:"last_name" json:"last_name" form:"last_name" binding:"required"`
	Email     string `db:"email" json:"email" form:"email" binding:"required"`
	Birthday  Date   `db:"birthday" json:"birthday" form:"birthday" binding:"required"`
}

type User struct {
	UserId
	NewUser
}
