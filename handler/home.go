package handler

import (
	"database/sql"
	"fitness-tracker/database"
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
	Label       string `json:"label"`
	Data        []int  `json:"data"`
	BorderWidth int    `json:"borderWidth"`
}

func HandleGetAppChartData(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID := c.Get("userID").(string)
		workoutID := c.PathParam("workoutID")
		startStr := c.QueryParam("start")
		start, err := time.Parse("2006-01-02 15:04:05", startStr)
		if err != nil {
			log.Println("invalid start:", err)
			return c.Render(http.StatusInternalServerError, "empty", nil)
		}
		end := time.Now().UTC()

		diff := end.Sub(start)
		days := int(diff.Hours() / 24.0)
		labels := make([]string, days)

		for start.Before(end) {
			labels = append(labels, start.Format("02/01"))
			start = start.Add(time.Duration(24 * time.Hour))
		}

		wes, err := database.ReadWorkoutEntriesBetweenTimes(db, userID, workoutID, start, end)

		// TODO: finish logic
	})
}
