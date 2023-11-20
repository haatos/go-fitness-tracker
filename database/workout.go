package database

import (
	"database/sql"
	"fitness-tracker/model"
	"fitness-tracker/schema"
	"time"

	"github.com/google/uuid"
)

func CreateWorkout(tx *sql.Tx, userID, workoutName string) (model.Workout, error) {
	w := model.Workout{}
	id := uuid.NewString()
	stmt, err := tx.Prepare(
		`
		INSERT INTO workout (id, name, user_id) values($1, $2, $3) RETURNING id, name, user_id
		`,
	)
	if err != nil {
		return w, err
	}
	err = stmt.QueryRow(id, workoutName, userID).Scan(&w.ID, &w.Name, &w.UserID)
	return w, err
}

func ReadAllWorkouts(db *sql.DB, userID string) ([]model.Workout, error) {
	workouts := []model.Workout{}
	stmt, err := db.Prepare(
		`
		SELECT id, name, user_id FROM workout
		where user_id = $1
		`,
	)
	if err != nil {
		return workouts, err
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return workouts, err
	}

	for rows.Next() {
		w := model.Workout{}

		if err := rows.Scan(&w.ID, &w.Name, &w.UserID); err != nil {
			return workouts, err
		}

		workouts = append(workouts, w)
	}

	return workouts, nil
}

func ReadWorkoutLastCreated(db *sql.DB, workoutID string) time.Time {
	var createdString string
	stmt, err := db.Prepare(
		`
		SELECT e.time
		FROM workout w
		INNER JOIN junction j
		ON j.workout_id = w.id
		INNER JOIN entry e
		ON e.junction_id = j.id
		WHERE w.id = $1
		ORDER BY e.time DESC
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
	created, err := time.Parse("2006-01-02T15:04:05Z", createdString)
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
		ON j.workout_id = w.id
		INNER JOIN exercise e
		ON e.id = j.exercise_id
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
