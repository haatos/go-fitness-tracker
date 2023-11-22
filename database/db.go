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

	return db
}
