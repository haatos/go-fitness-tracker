package database

import (
	"database/sql"
	"fitness-tracker/model"
)

func ReadAllWorkouts(db *sql.DB, userID string) ([]model.Workout, error) {
	stmt, err := db.Prepare(
		`
		SELECT id, name FROM workout
		where user_id = $1
		`,
	)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	workouts := []model.Workout{}
	for rows.Next() {
		w := model.Workout{}

		if err := rows.Scan(&w.ID, &w.Name); err != nil {
			return nil, err
		}

		workouts = append(workouts, w)
	}

	return workouts, nil
}
