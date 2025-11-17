package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/brandonrachal/gin-and-tonic/models"
	"github.com/brandonrachal/go-toolbox/dbutils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	DbConn             *sqlx.DB
	createUserStmt     *sqlx.Stmt
	getUserStmt        *sqlx.Stmt
	getFirstUserStmt   *sqlx.Stmt
	getUsersStmt       *sqlx.Stmt
	updateUserStmt     *sqlx.Stmt
	deleteUserStmt     *sqlx.Stmt
	deleteAllUsersStmt *sqlx.Stmt
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

	getUsersSql := `SELECT id, first_name, last_name, email, birthday FROM users`
	getUsersStmt, getUsersStmtErr := dbConn.Preparex(getUsersSql)
	if getUsersStmtErr != nil {
		return nil, getUsersStmtErr
	}

	getUserSql := fmt.Sprintf("%s where id = ?", getUsersSql)
	getUserStmt, getUserStmtErr := dbConn.Preparex(getUserSql)
	if getUserStmtErr != nil {
		return nil, getUserStmtErr
	}

	getFirstUserSql := fmt.Sprintf("%s limit 1", getUsersSql)
	getFirstUserStmt, getFirstUserStmtErr := dbConn.Preparex(getFirstUserSql)
	if getFirstUserStmtErr != nil {
		return nil, getFirstUserStmtErr
	}

	deleteUserSql := "delete from users where id = ?"
	deleteUserStmt, deleteUserStmtErr := dbConn.Preparex(deleteUserSql)
	if deleteUserStmtErr != nil {
		return nil, deleteUserStmtErr
	}
	deleteAllUsersSql := "delete from users"
	deleteAllUsersStmt, deleteAllUsersStmtErr := dbConn.Preparex(deleteAllUsersSql)
	if deleteAllUsersStmtErr != nil {
		return nil, deleteAllUsersStmtErr
	}
	updateUserSql := "update users set first_name = ?, last_name = ?, email = ?, birthday = ? where id = ?"
	updateUserStmt, updateUserStmtErr := dbConn.Preparex(updateUserSql)
	if updateUserStmtErr != nil {
		return nil, updateUserStmtErr
	}

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
		getFirstUserStmt:   getFirstUserStmt,
		updateUserStmt:     updateUserStmt,
		deleteUserStmt:     deleteUserStmt,
		deleteAllUsersStmt: deleteAllUsersStmt,
		getUsersStmt:       getUsersStmt,
		getUsersWithAgeSql: getUsersWithAgeSql,
		getAgeStatsSql:     getAgeStatsSql,
	}, nil
}

func (db *Client) CreateUser(ctx context.Context, firstName, lastName, email string, birthday time.Time) (sql.Result, error) {
	return db.createUserStmt.ExecContext(ctx, firstName, lastName, email, birthday)
}

func (db *Client) GetUser(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := db.getUserStmt.GetContext(ctx, &user, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *Client) GetFirstUser(ctx context.Context) (*models.User, error) {
	var user models.User
	err := db.getFirstUserStmt.GetContext(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *Client) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	rows, err := db.getUsersStmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Birthday)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (db *Client) UpdateUser(ctx context.Context, id int64, firstName, lastName, email string, birthday time.Time) (sql.Result, error) {
	return db.updateUserStmt.ExecContext(ctx, firstName, lastName, email, birthday, id)
}

func (db *Client) DeleteUser(ctx context.Context, id int64) (sql.Result, error) {
	return db.deleteUserStmt.ExecContext(ctx, id)
}

func (db *Client) DeleteAllUsers(ctx context.Context) (sql.Result, error) {
	return db.deleteAllUsersStmt.ExecContext(ctx)
}

func (db *Client) GetUsersWithAge(ctx context.Context) ([]models.UserWithAge, error) {
	var usersWithAage []models.UserWithAge
	err := db.DbConn.SelectContext(ctx, &usersWithAage, db.getUsersWithAgeSql)
	if err != nil {
		return nil, err
	}
	return usersWithAage, nil
}

func (db *Client) GetAgeStats(ctx context.Context) (*models.AgeStats, error) {
	var ageStats models.AgeStats
	err := db.DbConn.GetContext(ctx, &ageStats, db.getAgeStatsSql)
	if err != nil {
		return nil, err
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
	err = db.getUsersStmt.Close()
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
