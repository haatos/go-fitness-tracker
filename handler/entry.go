package handler

import (
	"database/sql"
	"encoding/json"
	"fitness-tracker/database"
	"fitness-tracker/model"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

func HandlePostWorkoutEntry(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID := c.Get("userID").(string)

		e := model.Entry{}
		if err := json.NewDecoder(c.Request().Body).Decode(&e); err != nil {
			log.Println("err decoding json:", err)
			c.Response().Header().Set("HX-Retarget", "#entry-error")
			return c.Render(http.StatusBadRequest, "error", struct {
				ID      string
				Message string
			}{
				ID:      "entry-error",
				Message: "Invalid data",
			})
		}

		e, err = database.CreateEntry(db, e)
	})
}
