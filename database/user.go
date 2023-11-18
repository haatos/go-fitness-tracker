package database

import (
	"database/sql"
	"fitness-tracker/model"

	"github.com/google/uuid"
)

func CreateUser(db *sql.DB, email, passwordHash string) error {
	user := model.User{}
	userID := uuid.NewString()
	stmt, err := db.Prepare(
		`
		INSERT INTO user (id, email, password_hash) values($1, $2, $3)
		RETURNING id, email, password_hash
		`,
	)
	if err != nil {
		return err
	}

	return stmt.QueryRow(userID, email, passwordHash).Scan(
		&user.ID, &user.Email, &user.PasswordHash,
	)
}

func ReadUserByEmail(db *sql.DB, email string) (model.User, error) {
	user := model.User{}
	stmt, err := db.Prepare(
		`
		SELECT id, email, password_hash, created_on FROM user WHERE email = $1
		`,
	)
	if err != nil {
		return user, err
	}

	err = stmt.QueryRow(email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	return user, err
}
