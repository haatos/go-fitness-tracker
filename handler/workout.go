package handler

import (
	"database/sql"
	"fitness-tracker/database"
	"fitness-tracker/model"
	"fitness-tracker/schema"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
)

func HandleGetWorkoutID(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		workoutID := c.PathParam("id")
		userID := c.Get("userID").(string)

		latestCreated := database.ReadWorkoutLastCreated(db, workoutID)

		entries := []model.Entry{}

		if !latestCreated.IsZero() {
			start := time.Date(latestCreated.Year(), latestCreated.Month(), latestCreated.Day(), 0, 0, 0, 0, latestCreated.Location())
			end := start.Add(24 * time.Hour)

			var err error
			entries, err = database.ReadEntriesBetweenTimes(db, userID, start, end)
			if err != nil {
				log.Println("err reading entries between times:", err)
			}
		}

		wos, err := database.ReadWorkoutJunctions(db, userID, workoutID)
		if err != nil {
			log.Println("err reading workout junctions:", err)
		}

		for i := range wos {
			sets := []schema.Set{}

			for j := 1; j <= wos[i].SetCount; j++ {
				s := schema.Set{}
				for _, entry := range entries {
					if entry.SetNumber == j && wos[i].JunctionID == entry.JunctionID {
						s.Weight = entry.Weight
						s.Reps = entry.Reps
						break
					}
				}
				s.JunctionID = wos[i].JunctionID
				s.SetNumber = j
				sets = append(sets, s)
			}
			wos[i].Sets = sets
		}

		return c.Render(http.StatusOK, "workout-id", struct {
			JunctionID string
			Data       []schema.WorkoutOut
		}{
			JunctionID: wos[0].JunctionID,
			Data:       wos,
		})
	})
}

func HandlePostWorkoutAddExercise(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		exerciseID := c.PathParam("id")
		userID := c.Get("userID").(string)

		ex, err := database.ReadExerciseByID(db, exerciseID)
		if ex.UserID != nil && *ex.UserID != userID {
			return c.Render(http.StatusUnauthorized, "empty", nil)
		}
		if err != nil {
			log.Println("err reading exercise by id:", err)
			return c.Render(http.StatusInternalServerError, "empty", nil)
		}

		return c.Render(http.StatusOK, "workout-exercise", ex)
	})
}

func HandleGetWorkoutCreate(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID := c.Get("userID").(string)
		exercises, err := database.ReadAllUserExercises(db, userID)
		if err != nil {
			log.Println("err reading user exercises:", err)
		}

		defaultExercises := []model.Exercise{}
		userExercises := []model.Exercise{}

		for _, e := range exercises {
			if e.UserID == nil {
				defaultExercises = append(defaultExercises, e)
			} else if *e.UserID == userID {
				userExercises = append(userExercises, e)
			}
		}

		return c.Render(
			http.StatusOK,
			"create-workout",
			struct {
				UserExercises    []model.Exercise
				DefaultExercises []model.Exercise
			}{
				UserExercises: userExercises, DefaultExercises: defaultExercises,
			})
	})
}

func HandlePostWorkoutCreate(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID := c.Get("userID").(string)
		workoutName := c.FormValue("workout_name")

		tx, err := db.Begin()
		if err != nil {
			log.Println("unable to begin tx:", err)
			return c.Render(http.StatusInternalServerError, "create-workout", nil)
		}
		defer tx.Rollback()

		w, err := database.CreateWorkout(tx, userID, workoutName)

		data := map[string][]string{}
		for k, v := range c.Request().PostForm {
			data[k] = v
		}

		exerciseIDs := data["exercise_id"]
		setCountStrings := data["set_count"]
		setCounts := []int{}
		for _, v := range setCountStrings {
			i, err := strconv.Atoi(v)
			if err != nil {
				return c.Render(http.StatusBadRequest, "create-workout", nil)
			}
			setCounts = append(setCounts, i)
		}

		if len(exerciseIDs) != len(setCounts) {
			log.Println("unequal amount of exercises and set counts")
			return c.Render(http.StatusBadRequest, "create-workout", nil)
		}

		for i := range exerciseIDs {
			id := exerciseIDs[i]
			setCount := setCounts[i]

			_, err := database.CreateJunction(tx, id, w.ID, userID, setCount)

			if err != nil {
				log.Println("err creating junction:", err)
				return c.Render(http.StatusInternalServerError, "create-workout", nil)
			}
		}

		err = tx.Commit()
		if err != nil {
			log.Println("err commiting tx:", err)
			return c.Render(http.StatusInternalServerError, "create-workout", nil)
		}

		return c.Redirect(http.StatusSeeOther, "/app")
	})
}
