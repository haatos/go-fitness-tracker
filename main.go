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
	session.Store = sessions.NewCookieStore([]byte(os.Getenv("FIT_SESSION_SECRET")))

	db := database.InitializeDatabase()

	e := echo.New()

	e.Renderer = &Templates{
		templates: parseTemplates(),
	}

	e.Static("/static", "static")

	publicGroup := e.Group("")
	publicGroup.Use(session.RedirectMiddleware)

	publicGroup.POST("/remove-me", func(c echo.Context) error {
		return c.Render(http.StatusOK, "empty", nil)
	})

	publicGroup.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "login", nil)
	})
	publicGroup.GET("/register", func(c echo.Context) error {
		return c.Render(http.StatusOK, "register", nil)
	})
	publicGroup.POST("/register", handler.HandlePostRegister(db))
	publicGroup.POST("/login", handler.HandleLogIn(db))

	appGroup := e.Group("/app")
	appGroup.Use(session.SessionMiddleware(db))
	appGroup.GET("", handler.HandleGetAppHome(db))
	appGroup.GET("/chart-data", handler.HandleGetAppChartData(db))

	exerciseGroup := appGroup.Group("/exercise")
	exerciseGroup.POST("/add", handler.HandlePostExerciseAdd(db))

	workoutGroup := appGroup.Group("/workout")
	workoutGroup.GET("/:id", handler.HandleGetWorkoutID(db))
	workoutGroup.POST("/add-exercise/:id", handler.HandlePostWorkoutAddExercise(db))
	workoutGroup.GET("/create", handler.HandleGetWorkoutCreate(db))
	workoutGroup.POST("/create", handler.HandlePostWorkoutCreate(db))
	workoutGroup.POST("/entry", handler.HandlePostWorkoutEntry(db))
	workoutGroup.PATCH("/entry/:id", handler.HandlePatchWorkoutEntryID(db))

	e.Start(os.Getenv("FIT_PORT"))
}
