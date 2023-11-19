package database

import (
	"database/sql"
	"fitness-tracker/model"
	"time"

	"github.com/google/uuid"
)

func CreateEntry(db *sql.DB, e model.Entry) error {
	id := uuid.NewString()
	stmt, err := db.Prepare(
		`
		INSERT INTO entry (id, user_id, junction_id, weight, reps, set_number, time) values(
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id
		`,
	)
	if err != nil {
		return err
	}
	return stmt.QueryRow(id, e.UserID, e.JunctionID, e.Weight, e.Reps, e.SetNumber, e.Time.Format("2006-01-02 15:04:05")).Scan(&e.ID)
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
