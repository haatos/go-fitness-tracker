package database

import (
	"database/sql"

	"github.com/tomihaapalainen/go-fitness-tracker/model"

	"github.com/google/uuid"
)

func CreateExercise(db *sql.DB, userID, name string) (model.Exercise, error) {
	e := model.Exercise{}
	id := uuid.NewString()
	stmt, err := db.Prepare(
		`
		INSERT INTO exercise (id, name, user_id) values($1, $2, $3) RETURNING id, name, user_id
		`,
	)
	if err != nil {
		return e, err
	}
	err = stmt.QueryRow(id, name, userID).Scan(&e.ID, &e.Name, &e.UserID)
	return e, err
}

func ReadExerciseByID(db *sql.DB, exerciseID string) (model.Exercise, error) {
	e := model.Exercise{}
	stmt, err := db.Prepare(
		`
		SELECT id, name, user_id FROM exercise
		WHERE id = $1
		`,
	)
	if err != nil {
		return e, err
	}

	err = stmt.QueryRow(exerciseID).Scan(&e.ID, &e.Name, &e.UserID)
	return e, err
}

func ReadAllUserExercises(db *sql.DB, userID string) ([]model.Exercise, error) {
	exs := []model.Exercise{}
	stmt, err := db.Prepare(
		`
		SELECT id, name, user_id FROM exercise
		WHERE user_id IS NULL OR user_id = $1
		`,
	)
	if err != nil {
		return exs, err
	}
	rows, err := stmt.Query(userID)
	if err != nil {
		return exs, err
	}
	for rows.Next() {
		ex := model.Exercise{}
		if err := rows.Scan(&ex.ID, &ex.Name, &ex.UserID); err != nil {
			return exs, err
		}

		exs = append(exs, ex)
	}

	return exs, nil
}
