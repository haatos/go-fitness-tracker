package handler

import (
	"database/sql"
	"fitness-tracker/database"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

func HandlePostExerciseAdd(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID := c.Get("userID").(string)
		name := c.FormValue("name")
		name = strings.TrimSpace(name)

		if name == "" {
			log.Println("empty exercise name given")
			return c.Render(http.StatusBadRequest, "empty", nil)
		}

		ex, err := database.CreateExercise(db, userID, name)
		if err != nil {
			log.Println("err creating exercise:", err)
			return c.Render(http.StatusInternalServerError, "empty", nil)
		}

		return c.Render(http.StatusOK, "add-exercise-button", ex)
	})
}
