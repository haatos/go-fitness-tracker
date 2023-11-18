package session

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
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
		c.Set("userID", userID)
		return next(c)
	}
}
