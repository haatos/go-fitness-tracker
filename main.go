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

	e.Static("/static", "static")

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

	appGroup := e.Group("/app")
	appGroup.Use(session.SessionMiddleware(db))
	appGroup.GET("", handler.HandleGetAppHome(db))
	appGroup.GET("/chart-data/:workoutID", handler.HandleGetAppChartData(db))

	exerciseGroup := appGroup.Group("/exercise")
	exerciseGroup.POST("/add", handler.HandlePostExerciseAdd(db))

	workoutGroup := appGroup.Group("/workout")
	workoutGroup.GET("/:id", handler.HandleGetWorkoutID(db))
	workoutGroup.POST("/add-exercise/:id", handler.HandlePostWorkoutAddExercise(db))
	workoutGroup.GET("/create", handler.HandleGetWorkoutCreate(db))
	workoutGroup.POST("/create", handler.HandlePostWorkoutCreate(db))
	workoutGroup.POST("/entry", handler.HandlePostWorkoutEntry(db))
	workoutGroup.PATCH("/entry/:id", handler.HandlePatchWorkoutEntryID(db))

	e.Start(":8080")
}
