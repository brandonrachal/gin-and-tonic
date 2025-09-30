package db_client

import (
	"context"
	"database/sql"

	"github.com/brandonrachal/go-toolbox/dbutils"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DBClient struct {
	DbConn             *sqlx.DB
	createUserTableSql string
	insertUserStmt     *sqlx.Stmt
	getUserStmt        *sqlx.Stmt
	updateUserStmt     *sqlx.Stmt
	getAllUsersSql     string
}

// ":memory:"

func NewDBClient(dataSourceName string) (*DBClient, error) {
	dbConn, dbConnErr := dbutils.NewSQLiteDBConn(dataSourceName)
	if dbConnErr != nil {
		return nil, dbConnErr
	}
	createUserTableSql := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
	);`
	insertUserStmt, insertUserStmtErr := dbConn.Preparex("insert into users (name, email) values (?, ?)")
	if insertUserStmtErr != nil {
		return nil, insertUserStmtErr
	}
	getUserStmt, getUserStmtErr := dbConn.Preparex("select id, name, email from users where id = ?")
	if getUserStmtErr != nil {
		return nil, getUserStmtErr
	}
	updateUserStmt, updateUserStmtErr := dbConn.Preparex("update users set name = ?, email = ? where id = ?")
	if updateUserStmtErr != nil {
		return nil, updateUserStmtErr
	}
	getAllUsersSql := `SELECT id, name, email FROM users;`

	return &DBClient{
		DbConn:             dbConn,
		createUserTableSql: createUserTableSql,
		insertUserStmt:     insertUserStmt,
		getUserStmt:        getUserStmt,
		updateUserStmt:     updateUserStmt,
		getAllUsersSql:     getAllUsersSql,
	}, nil
}

func (db *DBClient) CreateUserTable(ctx context.Context) (sql.Result, error) {
	return db.DbConn.ExecContext(ctx, db.createUserTableSql)
}

func (db *DBClient) InsertUser(ctx context.Context, name, email string) (sql.Result, error) {
	return db.insertUserStmt.ExecContext(ctx, name, email)
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

func (db *DBClient) UpdateUser(ctx context.Context, name, email string, id int) (sql.Result, error) {
	return db.updateUserStmt.ExecContext(ctx, name, email, id)
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
