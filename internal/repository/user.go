package repository

import (
	"database/sql"

	"birthday-bot/internal/model"
)

func GetUserByUserName(db *sql.DB, username string) (*model.User, error) {
	var user model.User
	err := db.QueryRow("SELECT username, password FROM users WHERE username=?",
		username).Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SaveUser(db *sql.DB, user model.User) error {
	_, err := db.Exec(`INSERT INTO users(username, password, role) VALUES(?, ?, ?)`,
		user.Username, user.Password, user.Role)
	return err
}
