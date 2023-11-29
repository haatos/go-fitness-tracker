package schema

import "time"

type Entry struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	JunctionID   string    `json:"junction_id"`
	SetNumber    int       `json:"set_number,string"`
	Weight       int       `json:"weight,string"`
	Reps         int       `json:"reps,string"`
	Time         time.Time `json:"time"`
	ExerciseName string    `json:"exercise_name"`
}
