package schema

import "time"

type Set struct {
	JunctionID string
	SetNumber  int
	Weight     int
	Reps       int
}

type WorkoutOut struct {
	WorkoutName  string
	ExerciseName string
	JunctionID   string
	SetCount     int
	Sets         []Set
}

type WorkoutEntry struct {
	ExerciseName string
	SetNumber    int
	Performance  int
	Time         time.Time
}
