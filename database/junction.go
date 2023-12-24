package database

import (
	"database/sql"

	"github.com/tomihaapalainen/go-fitness-tracker/model"

	"github.com/google/uuid"
)

func CreateJunction(tx *sql.Tx, exerciseID, workoutID, userID string, setCount int) (model.Junction, error) {
	id := uuid.NewString()
	j := model.Junction{}
	stmt, err := tx.Prepare(
		`
		INSERT INTO junction (id, exercise_id, workout_id, user_id, set_count) values(
			$1, $2, $3, $4, $5
		)
		RETURNING id, exercise_id, workout_id, user_id, set_count
		`,
	)
	if err != nil {
		return j, err
	}
	err = stmt.QueryRow(id, exerciseID, workoutID, userID, setCount).Scan(
		&j.ID, &j.ExerciseID, &j.WorkoutID, &j.UserID, &j.SetCount)

	return j, err
}
