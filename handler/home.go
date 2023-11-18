package handler

import (
	"database/sql"
	"fitness-tracker/database"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

func HandleGetAppHome(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID := c.Get("userID").(string)
		workouts, err := database.ReadAllWorkouts(db, userID)
		log.Println(workouts, err)
		return c.Render(http.StatusOK, "home", workouts)
	})
}
