package model

type Junction struct {
	ID         string `json:"id"`
	ExerciseID string `json:"exercise_id"`
	WorkoutID  string `json:"workout_id"`
	UserID     string `json:"user_id"`
	SetCount   int    `json:"set_count"`
}
