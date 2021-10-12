package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/jhallat/todo-schedule-service/cors"
	"github.com/jhallat/todo-schedule-service/httphelper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const scheduledTasksBasePath = "schedule"

func SetupRoutes(apiBasePath string) {
	handleTasks := http.HandlerFunc(tasksHandler)
	handleTask := http.HandlerFunc(taskHandler)
	handleRange := http.HandlerFunc(rangeHandler)
	http.Handle(fmt.Sprintf("%s/%s/range", apiBasePath, scheduledTasksBasePath), cors.Middleware(handleRange))
	http.Handle(fmt.Sprintf("%s/%s/range/", apiBasePath, scheduledTasksBasePath), cors.Middleware(handleRange))
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath,  scheduledTasksBasePath), cors.Middleware(handleTask))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath,  scheduledTasksBasePath), cors.Middleware(handleTasks))
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		httphelper.GetRequest(w, r, "schedule/:id", func(params map[string]string) (interface{}, error) {
			id := params["id"]
			taskId, err := strconv.Atoi(id)
			if err != nil { return nil, err }
			return getSchedule(taskId)
		})
	case http.MethodPut:
		var newSchedule Schedule
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newSchedule)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = saveSchedule(newSchedule)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		httphelper.GetQueryRequest(w, r,  func(values url.Values) (interface{}, error) {
			dayText := values.Get("day")
			day, err := time.Parse("2006-01-02", dayText)
			if err != nil { return nil, err }
			return getSchedulesForDate(day, false)
		})
	case http.MethodOptions:
		return
	}
}

func rangeHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("rageHandler called")
	switch r.Method {
	case http.MethodGet:
		httphelper.GetQueryRequest(w, r,  func(values url.Values) (interface{}, error) {
			startText := values.Get("start")
			start, err := time.Parse("2006-01-02", startText)
			if err != nil {
				log.Print(err)
				return nil, err
			}
			endText := values.Get("end")
			end, err := time.Parse("2006-01-02", endText)
			if err != nil {
				log.Print(err)
				return nil, err
			}
			return getScheduleForRange(start, end)
		})
	case http.MethodOptions:
		return
	}
}