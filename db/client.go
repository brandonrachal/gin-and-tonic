package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/brandonrachal/go-toolbox/dbutils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	DbConn             *sqlx.DB
	createUserStmt     *sqlx.Stmt
	getUserStmt        *sqlx.Stmt
	updateUserStmt     *sqlx.Stmt
	deleteUserStmt     *sqlx.Stmt
	getUsersSql        string
	getUsersWithAgeSql string
	getAgeStatsSql     string
}

func NewClient(dataSourceName string) (*Client, error) {
	dbConn, dbConnErr := dbutils.NewSQLiteDBConn(dataSourceName)
	if dbConnErr != nil {
		return nil, dbConnErr
	}
	createUserSql := "insert into users(first_name, last_name, email, birthday) values (?, ?, ?, ?)"
	createUserStmt, createUserStmtErr := dbConn.Preparex(createUserSql)
	if createUserStmtErr != nil {
		return nil, createUserStmtErr
	}
	getUserSql := "select id, first_name, last_name, email, birthday from users where id = ?"
	getUserStmt, getUserStmtErr := dbConn.Preparex(getUserSql)
	if getUserStmtErr != nil {
		return nil, getUserStmtErr
	}
	deleteUserSql := "delete from users where id = ?"
	deleteUserStmt, deleteUserStmtErr := dbConn.Preparex(deleteUserSql)
	if deleteUserStmtErr != nil {
		return nil, deleteUserStmtErr
	}
	updateUserSql := "update users set first_name = ?, last_name = ?, email = ?, birthday = ? where id = ?"
	updateUserStmt, updateUserStmtErr := dbConn.Preparex(updateUserSql)
	if updateUserStmtErr != nil {
		return nil, updateUserStmtErr
	}
	getUsersSql := `SELECT id, first_name, last_name, email, birthday FROM users`
	getUsersWithAgeSql := `select id, first_name, last_name, email, birthday, ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) AS age_in_years from users;`
	getAgeStatsSql := `select
    count(case when ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) < 13 then 1 end) as preteen,
    count(case when ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) > 12 and ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) < 20 then 1 end) as teens,
    count(case when ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) > 19 and ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) < 30 then 1 end) as twenties,
    count(case when ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) > 29 and ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) < 40 then 1 end) as thirties,
    count(case when ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) > 39 and ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) < 50 then 1 end) as forties,
    count(case when ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) > 49 and ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) < 60 then 1 end) as fifties,
    count(case when ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) > 59 and ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) < 70 then 1 end) as sixties,
    count(case when ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) > 69 and ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) < 80 then 1 end) as seventies,
    count(case when ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) > 79 and ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) < 90 then 1 end) as eighties,
    count(case when ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) > 89 and ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) < 100 then 1 end) as nineties,
    count(case when ROUND((JULIANDAY('now') - JULIANDAY(birthday)) / 365.25) > 99 then 1 end) as centurion
from
    users;`

	return &Client{
		DbConn:             dbConn,
		createUserStmt:     createUserStmt,
		getUserStmt:        getUserStmt,
		updateUserStmt:     updateUserStmt,
		deleteUserStmt:     deleteUserStmt,
		getUsersSql:        getUsersSql,
		getUsersWithAgeSql: getUsersWithAgeSql,
		getAgeStatsSql:     getAgeStatsSql,
	}, nil
}

func (db *Client) CreateUser(ctx context.Context, firstName, lastName, email string, birthday time.Time) (sql.Result, error) {
	return db.createUserStmt.ExecContext(ctx, firstName, lastName, email, birthday)
}

func (db *Client) User(ctx context.Context, id int) (*User, error) {
	var user User
	err := db.getUserStmt.GetContext(ctx, &user, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *Client) UpdateUser(ctx context.Context, id int, firstName, lastName, email string, birthday time.Time) (sql.Result, error) {
	return db.updateUserStmt.ExecContext(ctx, firstName, lastName, email, birthday, id)
}

func (db *Client) DeleteUser(ctx context.Context, id int) (sql.Result, error) {
	return db.deleteUserStmt.ExecContext(ctx, id)
}

func (db *Client) Users(ctx context.Context) ([]User, error) {
	var users []User
	err := db.DbConn.SelectContext(ctx, &users, db.getUsersSql)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (db *Client) UsersWithAge(ctx context.Context) ([]UserWithAge, error) {
	var usersWithAage []UserWithAge
	err := db.DbConn.SelectContext(ctx, &usersWithAage, db.getUsersWithAgeSql)
	if err != nil {
		return nil, err
	}
	return usersWithAage, nil
}

func (db *Client) AgeStats(ctx context.Context) (*AgeStats, error) {
	var ageStats AgeStats
	err := db.DbConn.GetContext(ctx, &ageStats, db.getAgeStatsSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &ageStats, nil
}

func (db *Client) Close() error {
	var err error
	err = db.createUserStmt.Close()
	if err != nil {
		return err
	}
	err = db.getUserStmt.Close()
	if err != nil {
		return err
	}
	err = db.updateUserStmt.Close()
	if err != nil {
		return err
	}
	err = db.deleteUserStmt.Close()
	if err != nil {
		return err
	}
	err = db.DbConn.Close()
	if err != nil {
		return err
	}
	return nil
}
