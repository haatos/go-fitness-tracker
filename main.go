package main

import (
	"fitness-tracker/database"
	"fitness-tracker/handler"
	"fitness-tracker/session"
	"net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	db := database.InitializeDatabase()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})
	e.GET("/register", func(c echo.Context) error {
		return c.Render(http.StatusOK, "register", nil)
	})
	e.POST("/register", handler.HandlePostRegister(db))
	e.POST("/login", handler.HandleLogIn(db))

	g := e.Group("/app")
	g.Use(session.SessionMiddleware)

	g.GET("", handler.HandleGetAppHome(db))

	g.GET("/workout/:id", handler.HandleGetWorkoutID(db))
}
