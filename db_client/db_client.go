package db_client

import (
	"context"
	"database/sql"

	"github.com/brandonrachal/go-toolbox/dbutils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DBClient struct {
	DbConn         *sqlx.DB
	insertUserStmt *sqlx.Stmt
	getUserStmt    *sqlx.Stmt
	updateUserStmt *sqlx.Stmt
	getAllUsersSql string
}

func NewDBClient(dataSourceName string) (*DBClient, error) {
	dbConn, dbConnErr := dbutils.NewSQLiteDBConn(dataSourceName)
	if dbConnErr != nil {
		return nil, dbConnErr
	}
	insertUserSql := "insert into users (first_name, last_name, email) values (?, ?, ?)"
	insertUserStmt, insertUserStmtErr := dbConn.Preparex(insertUserSql)
	if insertUserStmtErr != nil {
		return nil, insertUserStmtErr
	}
	getUserSql := "select id, first_name, last_name, email from users where id = ?"
	getUserStmt, getUserStmtErr := dbConn.Preparex(getUserSql)
	if getUserStmtErr != nil {
		return nil, getUserStmtErr
	}
	updateUserSql := "update users set first_name = ?, last_name = ?, email = ? where id = ?"
	updateUserStmt, updateUserStmtErr := dbConn.Preparex(updateUserSql)
	if updateUserStmtErr != nil {
		return nil, updateUserStmtErr
	}
	getAllUsersSql := `SELECT id, first_name, last_name, email FROM users;`
	return &DBClient{
		DbConn:         dbConn,
		insertUserStmt: insertUserStmt,
		getUserStmt:    getUserStmt,
		updateUserStmt: updateUserStmt,
		getAllUsersSql: getAllUsersSql,
	}, nil
}

func (db *DBClient) InsertUser(ctx context.Context, firstName, lastName, email string) (sql.Result, error) {
	return db.insertUserStmt.ExecContext(ctx, firstName, lastName, email)
}

func (db *DBClient) GetUser(ctx context.Context, id int) (*User, error) {
	var user User
	err := db.getUserStmt.GetContext(ctx, &user, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DBClient) GetAllUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := db.DbConn.SelectContext(ctx, &users, db.getAllUsersSql)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (db *DBClient) UpdateUser(ctx context.Context, firstName, lastName, email string, id int) (sql.Result, error) {
	return db.updateUserStmt.ExecContext(ctx, firstName, lastName, email, id)
}

func (db *DBClient) Close() error {
	var err error
	err = db.insertUserStmt.Close()
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
	err = db.DbConn.Close()
	if err != nil {
		return err
	}
	return nil
}
