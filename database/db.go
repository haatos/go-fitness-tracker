package database

import (
	"database/sql"
	"log"
)

func InitializeDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "file:.///db.sqlite3?_fk=ON")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		`
		CREATE TABLE IF NOT EXISTS user (
			id TEXT PRIMARY KEY,
			email TEXT,
			password_hash TEXT,
			created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
		`,
	)

	_, err = db.Exec(
		`
		CREATE TABLE IF NOT EXISTS exercise (
			id TEXT PRIMARY KEY,
			name TEXT,
			user_id TEXT,
			FOREIGN KEY(user_id) REFERENCES user(id)
		)
		`,
	)

	_, err = db.Exec(
		`
		CREATE TABLE IF NOT EXISTS workout (
			id TEXT PRIMARY KEY,
			name TEXT,
			user_id TEXT,
			FOREIGN KEY(user_id) REFERENCES user(id)
		)
		`,
	)

	_, err = db.Exec(
		`
		CREATE TABLE IF NOT EXISTS junction (
			exercise_id TEXT,
			workout_id TEXT,
			user_id TEXT,
			set_count INTEGER,
			FOREIGN KEY(exercise_id) REFERENCES exercise(id),
			FOREIGN KEY(workout_id) REFERENCES workout(id),
			FOREIGN KEY(user_id) REFERENCES user(id)
		)
		`,
	)

	_, err = db.Exec(
		`
		CREATE TABLE IF NOT EXISTS entry (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			junction_id TEXT,
			set_number INTEGER,
			weight INTEGER,
			reps INTEGER,
			FOREIGN KEY(user_id) REFERENCES user(id),
			FOREIGN KEY(junction_id) REFERENCES junction(id)
		)
		`,
	)

	return db
}
