package handler

import (
	"database/sql"
	"encoding/json"
	"fitness-tracker/database"
	"fitness-tracker/model"
	"log"
	"net/http"
	"time"

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
		e.UserID = userID
		e.Time = time.Now().UTC()

		e, err := database.CreateEntry(db, e)
		if err != nil {
			log.Println("err creating entry:", err)
			c.Response().Header().Set("HX-Retarget", "#entry-error")
			return c.Render(http.StatusBadRequest, "error", struct {
				ID      string
				Message string
			}{
				ID:      "entry-error",
				Message: "Unable to create entry, try again later",
			})
		}

		return c.Render(http.StatusOK, "update-entry", e)
	})
}

func HandlePatchWorkoutEntryID(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID := c.Get("userID").(string)
		id := c.PathParam("id")

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
		e.ID = id
		e.UserID = userID

		e, err := database.PatchEntry(db, e)
		if err != nil {
			log.Println("err patching entry:", err)
			c.Response().Header().Set("HX-Retarget", "#entry-error")
			return c.Render(http.StatusBadRequest, "error", struct {
				ID      string
				Message string
			}{
				ID:      "entry-error",
				Message: "Unable to update entry, try again later",
			})
		}

		return c.Render(http.StatusOK, "update-entry", e)
	})
}
