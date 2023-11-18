package database

import (
	"database/sql"
	"fitness-tracker/schema"
	"time"
)

func ReadWorkoutLastCreated(db *sql.DB, workoutID string) time.Time {
	var createdString string
	stmt, err := db.Prepare(
		`
		SELECT e.created
		FROM workout w
		INNER JOIN junction j
		ON j.workout = w.id
		INNER JOIN entry e
		ON e.junction = j.id
		WHERE w.id = $1
		LIMIT 1
		`,
	)
	if err != nil {
		return time.Time{}
	}

	err = stmt.QueryRow(workoutID).Scan(&createdString)
	if err != nil {
		return time.Time{}
	}
	created, err := time.Parse("2006-01-02 15:04:05Z", createdString)
	if err != nil {
		return time.Time{}
	}

	return created
}

func ReadWorkoutJunctions(db *sql.DB, userID, workoutID string) ([]schema.WorkoutOut, error) {
	stmt, err := db.Prepare(
		`
		SELECT e.name, w.name, j.id, j.set_count
		FROM workout w
		INNER JOIN junction j
		ON j.workout = w.id
		INNER JOIN exercise e
		ON e.id = j.exercise
		WHERE w.user_id = $1 AND w.id = $2
		`,
	)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userID, workoutID)
	if err != nil {
		return nil, err
	}

	wos := []schema.WorkoutOut{}
	for rows.Next() {
		wo := schema.WorkoutOut{}

		if err := rows.Scan(&wo.ExerciseName, &wo.WorkoutName, &wo.JunctionID, &wo.SetCount); err != nil {
			return nil, err
		}

		wos = append(wos, wo)
	}

	return wos, nil
}
