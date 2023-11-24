package handler

import (
	"database/sql"
	"fitness-tracker/database"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

func HandleGetAppHome(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID := c.Get("userID").(string)
		workouts, err := database.ReadAllWorkouts(db, userID)
		if err != nil {
			log.Println("err reading all workouts:", err)
		}

		return c.Render(http.StatusOK, "home", workouts)
	})
}

type Dataset struct {
	Label string `json:"label"`
	Data  []int  `json:"data"`
}

type OutputData struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

type Data struct {
	Data []OutputData `json:"data"`
}

func HandleGetAppChartData(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID := c.Get("userID").(string)
		workoutID := c.PathParam("workoutID")
		startStr := c.QueryParam("start")
		log.Println(startStr)
		start, err := time.Parse("2006-01-02 15:04:05", startStr)
		if err != nil {
			log.Println("invalid start:", err)
			return c.Render(http.StatusInternalServerError, "empty", nil)
		}
		end := time.Now().UTC()

		wes, err := database.ReadWorkoutEntriesBetweenTimes(db, userID, workoutID, start, end)

		labels := map[string]map[time.Time]bool{}
		output := map[string]map[int]*Dataset{}

		for _, we := range wes {
			if _, ok := output[we.ExerciseName]; ok {
				if _, ok := output[we.ExerciseName][we.SetNumber]; ok {
					output[we.ExerciseName][we.SetNumber].Data = append(output[we.ExerciseName][we.SetNumber].Data, we.Performance)
				} else {
					output[we.ExerciseName][we.SetNumber] = &Dataset{
						Label: fmt.Sprintf("Set %d", we.SetNumber),
						Data:  []int{we.Performance},
					}
				}
			} else {
				output[we.ExerciseName] = map[int]*Dataset{}
				output[we.ExerciseName][we.SetNumber] = &Dataset{
					Label: fmt.Sprintf("Set %d", we.SetNumber),
					Data:  []int{we.Performance},
				}
			}
			if _, ok := labels[we.ExerciseName]; ok {
				labels[we.ExerciseName][we.Time] = true
			} else {
				labels[we.ExerciseName] = map[time.Time]bool{}
				labels[we.ExerciseName][we.Time] = true
			}
		}

		// TODO add LABELS

		return c.JSON(http.StatusOK, output)
	})
}
