package handler

import (
	"database/sql"
	"fitness-tracker/database"
	"fitness-tracker/model"
	"fitness-tracker/schema"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
)

func HandleGetAppHome(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID := c.Get("userID").(string)
		workouts, err := database.ReadAllWorkouts(db, userID)
		if err != nil {
			log.Println("err reading all workouts:", err)
		}

		name, date, err := database.ReadLastWorkout(db, userID)
		if err != nil {
			log.Println("err reading last workout:", err)
		}

		data := struct {
			LastWorkout schema.LastWorkout
			Workouts    []model.Workout
			Env         string
		}{
			LastWorkout: schema.LastWorkout{
				Name: name,
				Time: date,
			},
			Workouts: workouts,
			Env:      os.Getenv("FIT_ENVIRONMENT"),
		}

		log.Printf("%+v", data)

		return c.Render(http.StatusOK, "app", data)
	})
}

type Dataset struct {
	Label string `json:"label"`
	Data  []int  `json:"data"`
}

type OutputData struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

type Data struct {
	Data []OutputData `json:"data"`
}

func HandleGetAppChartData(db *sql.DB) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		userID := c.Get("userID").(string)
		startStr := c.QueryParam("start")
		start, err := time.Parse("2006-01-02T15:04:05.000Z", startStr)
		if err != nil {
			log.Println("invalid start:", err)
			return c.Render(http.StatusInternalServerError, "empty", nil)
		}
		end := time.Now().UTC()

		wes, err := database.ReadWorkoutEntriesBetweenTimes(db, userID, start, end)

		labels := map[string]map[string]bool{}
		output := map[string]map[string]*Dataset{}

		for _, we := range wes {
			set := fmt.Sprintf("Set %d", we.SetNumber)
			if _, ok := output[we.ExerciseName]; ok {
				if _, ok := output[we.ExerciseName][set]; ok {
					output[we.ExerciseName][set].Data = append(output[we.ExerciseName][set].Data, we.Performance)
				} else {
					output[we.ExerciseName][set] = &Dataset{
						Label: set,
						Data:  []int{we.Performance},
					}
				}
			} else {
				output[we.ExerciseName] = map[string]*Dataset{}
				output[we.ExerciseName][set] = &Dataset{
					Label: set,
					Data:  []int{we.Performance},
				}
			}
			if _, ok := labels[we.ExerciseName]; ok {
				labels[we.ExerciseName][we.Time.Format("2006-01-02")] = true
			} else {
				labels[we.ExerciseName] = map[string]bool{}
				labels[we.ExerciseName][we.Time.Format("2006-01-02")] = true
			}
		}

		// TODO add LABELS

		outputLabels := map[string][]string{}
		for ex := range labels {
			outputLabels[ex] = []string{}
			ls := []string{}
			for l := range labels[ex] {
				ls = append(ls, l[5:])
			}
			slices.Sort(ls)
			for _, l := range ls {
				split := strings.Split(l, "-")
				d := split[1]
				if d == "1" {
					outputLabels[ex] = append(outputLabels[ex], l)
				} else {
					outputLabels[ex] = append(outputLabels[ex], d)
				}

			}
		}

		finalOutput := struct {
			Labels map[string][]string            `json:"labels"`
			Data   map[string]map[string]*Dataset `json:"data"`
		}{
			Labels: outputLabels,
			Data:   output,
		}

		return c.JSON(http.StatusOK, finalOutput)
	})
}
