package handler

import (
	"database/sql"
	"fitness-tracker/database"
	"fitness-tracker/session"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogIn(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		email := c.FormValue("identity")
		password := c.FormValue("password")

		user, err := database.ReadUserByEmail(db, email)

		if err != nil {
			log.Println("err reading user by email:", err)
			return c.Render(
				http.StatusSeeOther,
				"index",
				struct{ Error string }{
					Error: "Invalid email or password",
				},
			)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			log.Println("err comparing password:", err)
			return c.Render(
				http.StatusSeeOther,
				"index",
				struct{ Error string }{
					Error: "Invalid email or password",
				},
			)
		}

		sess, err := session.Store.New(c.Request(), "session")
		if err != nil {
			return c.Render(http.StatusSeeOther, "index", nil)
		}

		sess.Values["userID"] = user.ID
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			log.Println("err saving session:", err)
			return c.Render(http.StatusSeeOther, "index", nil)
		}

		return c.Redirect(http.StatusSeeOther, "/app")
	})
}
