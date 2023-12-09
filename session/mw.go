package session

import (
	"database/sql"
	"fitness-tracker/database"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

func RedirectMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		sess, _ := Store.Get(c.Request(), "session")
		_, ok := sess.Values["userID"].(string)
		if ok && c.Request().Method == http.MethodGet {
			return c.Redirect(http.StatusSeeOther, "/app")
		}
		return next(c)
	})
}

func SessionMiddleware(db *sql.DB) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			sess, err := Store.Get(c.Request(), "session")
			if err != nil {
				return c.Redirect(http.StatusSeeOther, "/")
			}
			userID, ok := sess.Values["userID"].(string)
			if !ok {
				sess.Options.MaxAge = -1
				err := sess.Save(c.Request(), c.Response())
				if err != nil {
					log.Println("err deleting cookie", err)
				}
				return c.Redirect(http.StatusSeeOther, "/")
			}
			_, err = database.ReadUserByID(db, userID)
			if err != nil {
				sess.Options.MaxAge = -1
				err := sess.Save(c.Request(), c.Response())
				if err != nil {
					log.Println("err deleting cookie", err)
				}
				return c.Redirect(http.StatusSeeOther, "/")
			}
			c.Set("userID", userID)
			return next(c)
		})
	}
}
