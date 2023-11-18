package database

import (
	"database/sql"
	"fitness-tracker/model"
	"time"
)

func ReadEntriesBetweenTimes(db *sql.DB, userID string, start, end time.Time) ([]model.Entry, error) {
	stmt, err := db.Prepare(
		`
		SELECT id, user_id, junction, weight, reps, set_number
		FROM entry
		WHERE user_id = $1 AND created BETWEEN $2 and $3
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
		if err := rows.Scan(&e.ID, &e.UserID, &e.JunctionID, &e.Weight, &e.Reps, &e.SetNumber); err != nil {
			return nil, err
		}

		entries = append(entries, e)
	}

	return entries, nil
}
