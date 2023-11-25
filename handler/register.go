package handler

import (
	"database/sql"
	"fitness-tracker/database"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

func HandlePostRegister(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		email := c.FormValue("email")
		password := c.FormValue("password")
		passwordVerification := c.FormValue("password_verification")

		if password != passwordVerification {
			return c.Render(http.StatusSeeOther, "register", nil)
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("FIT_SUGAR")+password), 4)
		if err != nil {
			log.Println("err hashing password:", err)
			return c.Render(http.StatusSeeOther, "register", struct{ Error string }{Error: "Error hashing your password"})
		}

		userCount, err := database.CountUsers(db)
		if err != nil {
			log.Println("err counting users", err)
			return c.Render(http.StatusSeeOther, "register", struct{ Error string }{Error: "Error counting users"})
		}

		if userCount >= 100 {
			log.Printf("max users hit: %d\n", userCount)
			return c.Render(http.StatusSeeOther, "register", struct{ Error string }{Error: "Maximum amount of users reached."})
		}

		if err := database.CreateUser(db, email, string(passwordHash)); err != nil {
			log.Println("err creating user:", err)
			return c.Render(
				http.StatusSeeOther,
				"register",
				struct{ Error string }{
					Error: fmt.Sprintf(
						"Error creating new user. User with email '%s' already exists.", email,
					),
				},
			)
		}

		return c.Redirect(http.StatusSeeOther, "/")
	})
}
