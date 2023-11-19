package model

import "time"

type Entry struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	JunctionID string    `json:"junction_id"`
	SetNumber  int       `json:"set_number"`
	Weight     int       `json:"weight"`
	Reps       int       `json:"reps"`
	Time       time.Time `json:"time"`
}
