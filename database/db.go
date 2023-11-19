package database

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
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
			FOREIGN KEY(user_id) REFERENCES user(id),
			UNIQUE(name, user_id)
		)
		`,
	)

	exercises := []string{
		"backsquat",
		"frontsquat",
		"deadlift",
		"romanian deadlift",
		"split squat",
		"bench press",
		"dumbbell bench press",
		"incline bench press",
		"incline dumbbell bench press",
		"pull up",
		"chin up",
		"cable pulldown",
		"cable row",
	}

	var backSquatID string
	db.QueryRow("SELECT id FROM exercise WHERE name = $1", "backsquat").Scan(&backSquatID)

	if backSquatID == "" {
		for _, e := range exercises {
			id := uuid.NewString()
			_, err := db.Exec(
				`
			INSERT INTO exercise (id, name) values($1, $2)
			`,
				id, e,
			)
			if err != nil {
				log.Println("err adding default exercise:", err)
			}
		}
	}

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
			id TEXT PRIMARY KEY,
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
			time TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES user(id),
			FOREIGN KEY(junction_id) REFERENCES junction(id)
		)
		`,
	)

	return db
}
