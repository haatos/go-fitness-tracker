package main

import (
	"fitness-tracker/database"
	"fitness-tracker/dotenv"
	"fitness-tracker/handler"
	"fitness-tracker/session"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"

	_ "github.com/mattn/go-sqlite3"
)

var pathRe *regexp.Regexp = regexp.MustCompile(`^.+\.html$`)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, e echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func parseTemplates() *template.Template {
	paths := []string{}

	filepath.Walk("templates", func(path string, info fs.FileInfo, err error) error {
		if pathRe.Match([]byte(path)) {
			paths = append(paths, path)
		}
		return nil
	})

	return template.Must(template.ParseFiles(paths...))
}

func main() {
	dotenv.ParseDotenv()
	session.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	db := database.InitializeDatabase()

	e := echo.New()
	e.Renderer = &Templates{
		templates: parseTemplates(),
	}

	e.POST("/remove-me", func(c echo.Context) error {
		return c.Render(http.StatusOK, "empty", nil)
	})

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})
	e.GET("/register", func(c echo.Context) error {
		return c.Render(http.StatusOK, "register", nil)
	})
	e.POST("/register", handler.HandlePostRegister(db))
	e.POST("/login", handler.HandleLogIn(db))

	g := e.Group("/app")
	g.Use(session.SessionMiddleware(db))
	g.GET("", handler.HandleGetAppHome(db))
	g.POST("/exercise/add", handler.HandlePostExerciseAdd(db))
	g.GET("/workout/:id", handler.HandleGetWorkoutID(db))
	g.POST("/workout/add-exercise/:id", handler.HandlePostWorkoutAddExercise(db))
	g.GET("/workout/create", handler.HandleGetWorkoutCreate(db))
	g.POST("/workout/create", handler.HandlePostWorkoutCreate(db))
	g.POST("/workout/entry", handler.HandlePostWorkoutEntry(db))
	// g.PUT("/workout/entry/:id", handler.HandlePatchWorkoutEntryID(db))

	e.Start(":8080")
}
