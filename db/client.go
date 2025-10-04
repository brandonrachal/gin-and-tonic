package db

import (
	"context"
	"database/sql"

	"github.com/brandonrachal/go-toolbox/dbutils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	DbConn         *sqlx.DB
	createUserStmt *sqlx.Stmt
	getUserStmt    *sqlx.Stmt
	updateUserStmt *sqlx.Stmt
	deleteUserStmt *sqlx.Stmt
	getAllUsersSql string
}

func NewClient(dataSourceName string) (*Client, error) {
	dbConn, dbConnErr := dbutils.NewSQLiteDBConn(dataSourceName)
	if dbConnErr != nil {
		return nil, dbConnErr
	}
	insertUserSql := "insert into users(first_name, last_name, email) values (?, ?, ?)"
	createUserStmt, createUserStmtErr := dbConn.Preparex(insertUserSql)
	if createUserStmtErr != nil {
		return nil, createUserStmtErr
	}
	getUserSql := "select id, first_name, last_name, email from users where id = ?"
	getUserStmt, getUserStmtErr := dbConn.Preparex(getUserSql)
	if getUserStmtErr != nil {
		return nil, getUserStmtErr
	}
	deleteUserSql := "delete from users where id = ?"
	deleteUserStmt, deleteUserStmtErr := dbConn.Preparex(deleteUserSql)
	if deleteUserStmtErr != nil {
		return nil, deleteUserStmtErr
	}
	updateUserSql := "update users set first_name = ?, last_name = ?, email = ? where id = ?"
	updateUserStmt, updateUserStmtErr := dbConn.Preparex(updateUserSql)
	if updateUserStmtErr != nil {
		return nil, updateUserStmtErr
	}
	getAllUsersSql := `SELECT id, first_name, last_name, email FROM users`
	return &Client{
		DbConn:         dbConn,
		createUserStmt: createUserStmt,
		getUserStmt:    getUserStmt,
		updateUserStmt: updateUserStmt,
		deleteUserStmt: deleteUserStmt,
		getAllUsersSql: getAllUsersSql,
	}, nil
}

func (db *Client) CreateUser(ctx context.Context, firstName, lastName, email string) (sql.Result, error) {
	return db.createUserStmt.ExecContext(ctx, firstName, lastName, email)
}

func (db *Client) User(ctx context.Context, id int) (*User, error) {
	var user User
	err := db.getUserStmt.GetContext(ctx, &user, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *Client) UpdateUser(ctx context.Context, id int, firstName, lastName, email string) (sql.Result, error) {
	return db.updateUserStmt.ExecContext(ctx, firstName, lastName, email, id)
}

func (db *Client) DeleteUser(ctx context.Context, id int) (sql.Result, error) {
	return db.deleteUserStmt.ExecContext(ctx, id)
}

func (db *Client) AllUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := db.DbConn.SelectContext(ctx, &users, db.getAllUsersSql)
	if err != nil {
		return nil, err
	}
	return users, nil
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
