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

func ReadUserByID(db *sql.DB, id string) (model.User, error) {
	u := model.User{}
	stmt, err := db.Prepare(
		`
		SELECT id, email, password_hash, created_on FROM user WHERE id = $1
		`,
	)
	if err != nil {
		return u, err
	}
	err = stmt.QueryRow(id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedOn)
	return u, err
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

	err = stmt.QueryRow(email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedOn)
	return user, err
}

func CountUsers(db *sql.DB) (int, error) {
	stmt, err := db.Prepare(
		`
		SELECT COUNT(*) FROM user
		`,
	)
	if err != nil {
		return 0, err
	}
	var count int
	err = stmt.QueryRow().Scan(&count)
	return count, err
}
