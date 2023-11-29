package database

import (
	"database/sql"
	"fitness-tracker/model"
	"fitness-tracker/schema"
	"time"

	"github.com/google/uuid"
)

func CreateEntry(db *sql.DB, e model.Entry) (model.Entry, error) {
	id := uuid.NewString()
	stmt, err := db.Prepare(
		`
		INSERT INTO entry (id, user_id, junction_id, weight, reps, set_number, time) values(
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id
		`,
	)
	if err != nil {
		return e, err
	}
	err = stmt.QueryRow(
		id,
		e.UserID,
		e.JunctionID,
		e.Weight,
		e.Reps,
		e.SetNumber,
		e.Time.Format("2006-01-02 15:04:05"),
	).Scan(&e.ID)
	return e, err
}

func PatchEntry(db *sql.DB, e model.Entry) (model.Entry, error) {
	stmt, err := db.Prepare(
		`
		UPDATE entry
		SET weight = $1,
		    reps = $2
		WHERE id = $3
		RETURNING weight, reps
		`,
	)
	if err != nil {
		return e, err
	}
	err = stmt.QueryRow(e.Weight, e.Reps, e.ID).Scan(&e.Weight, &e.Reps)
	return e, err
}

func ReadEntriesBetweenTimes(db *sql.DB, userID string, start, end time.Time) ([]model.Entry, error) {
	stmt, err := db.Prepare(
		`
		SELECT id, user_id, junction_id, weight, reps, set_number, time
		FROM entry
		WHERE user_id = $1 AND time BETWEEN $2 and $3
		`,
	)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(
		userID, start.Format("2006-01-02 15:04:05Z"), end.Format("2006-01-02 15:04:05Z"))
	if err != nil {
		return nil, err
	}

	entries := []model.Entry{}
	for rows.Next() {
		e := model.Entry{}
		if err := rows.Scan(&e.ID, &e.UserID, &e.JunctionID, &e.Weight, &e.Reps, &e.SetNumber, &e.Time); err != nil {
			return nil, err
		}

		entries = append(entries, e)
	}

	return entries, nil
}

func ReadWorkoutEntriesBetweenTimes(db *sql.DB, userID string, start, end time.Time) ([]schema.WorkoutEntry, error) {
	wes := []schema.WorkoutEntry{}
	stmt, err := db.Prepare(
		`
		SELECT ex.name, en.set_number, en.weight, en.reps, en.time
		FROM entry en
		INNER JOIN junction j
		ON en.junction_id = j.id
		INNER JOIN exercise ex
		ON ex.id = j.exercise_id
		WHERE en.user_id = $1 AND en.time BETWEEN $2 AND $3
		`,
	)
	if err != nil {
		return wes, err
	}

	rows, err := stmt.Query(userID, start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	if err != nil {
		return wes, err
	}

	for rows.Next() {
		we := schema.WorkoutEntry{}

		var w int
		var r int
		if err := rows.Scan(&we.ExerciseName, &we.SetNumber, &w, &r, &we.Time); err != nil {
			return wes, err
		}

		if w == 0 {
			// non-weighted exercise
			we.Performance = r
		} else {
			// weighted exercise
			we.Performance = w * r
		}

		wes = append(wes, we)
	}

	return wes, nil
}

func ReadLatestEntryForEachExercise(db *sql.DB, userID string) ([]schema.Entry, error) {
	entries := []schema.Entry{}

	stmt, err := db.Prepare(
		`
		SELECT en.id, en.user_id, en.junction_id, en.set_number, en.weight, en.reps, en.time, ex.name, MAX(en.time)
		FROM entry en
		INNER JOIN junction j
		ON en.junction_id = j.id
		INNER JOIN exercise ex
		ON j.exercise_id = ex.id
		WHERE en.user_id = $1
		GROUP BY en.set_number, ex.name
		`,
	)
	if err != nil {
		return entries, err
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return entries, err
	}

	for rows.Next() {
		entry := schema.Entry{}

		var x interface{}
		if err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.JunctionID,
			&entry.SetNumber,
			&entry.Weight,
			&entry.Reps,
			&entry.Time,
			&entry.ExerciseName,
			&x,
		); err != nil {
			return entries, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
