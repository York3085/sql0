package dao

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

func AddUser(username, password string, db *sqlx.DB) error {
	// 检查用户名是否已存在
	user := User{}
	query := "SELECT username FROM users WHERE username = ?"
	err := db.Get(&user, query, username)
	if err == nil {
		return errors.New("user already exists")
	}

	// 插入用户信息
	query = "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err = db.Exec(query, username, password)
	if err != nil {
		return err
	}

	return nil
}

func SelectUser(username string, db *sqlx.DB) bool {
	user := User{}
	query := "SELECT username FROM users WHERE username = ?"
	err := db.Get(&user, query, username)
	if err == nil {
		return true
	}
	return false
}

func SelectPasswordFromUsername(username string, db *sqlx.DB) string {
	user := User{}
	query := "SELECT password FROM users WHERE username = ?"
	err := db.Get(&user, query, username)
	if err != nil {
		return ""
	}
	return user.Password
}
