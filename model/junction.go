package model

type Junction struct {
	ExerciseID string `json:"exercise_id"`
	WorkoutID  string `json:"workout_id"`
	UserID     string `json:"user_id"`
	SetCount   int    `json:"set_count"`
}
