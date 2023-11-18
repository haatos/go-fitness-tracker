package handler

import (
	"database/sql"
	"fitness-tracker/database"
	"fitness-tracker/model"
	"fitness-tracker/schema"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

func HandleGetWorkoutID(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		workoutID := c.PathParam("id")
		userID := c.Get("userID").(string)

		latestCreated := database.ReadWorkoutLastCreated(db, workoutID)

		entries := []model.Entry{}

		if !latestCreated.IsZero() {
			start := time.Date(latestCreated.Year(), latestCreated.Month(), latestCreated.Day(), 0, 0, 0, 0, latestCreated.Location())
			end := start.Add(24 * time.Hour)

			var err error
			entries, err = database.ReadEntriesBetweenTimes(db, userID, start, end)
			if err != nil {
				log.Println("err reading entries between times:", err)
			}
		}

		wos, err := database.ReadWorkoutJunctions(db, userID, workoutID)
		if err != nil {
			log.Println("err reading workout junctions:", err)
		}

		for i := range wos {
			sets := []schema.Set{}

			for j := 1; j <= wos[i].SetCount; j++ {
				s := schema.Set{}
				for _, entry := range entries {
					if entry.SetNumber == j && wos[i].JunctionID == entry.JunctionID {
						s.Weight = entry.Weight
						s.Reps = entry.Reps
						break
					}
				}
				s.JunctionID = wos[i].JunctionID
				s.SetNumber = j
				sets = append(sets, s)
			}
			wos[i].Sets = sets
		}

		return c.Render(http.StatusOK, "workout-id", struct {
			JunctionID string
			Data       []schema.WorkoutOut
		}{
			JunctionID: wos[0].JunctionID,
			Data:       wos,
		})
	})
}
