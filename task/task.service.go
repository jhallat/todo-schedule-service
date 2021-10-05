package task

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

const scheduledTasksBasePath = "scheduled-task"

func SetupRoutes(apiBasePath string) {
	handleTasksV2 := http.HandlerFunc(tasksHandlerV2)
	handleTaskV2 := http.HandlerFunc(taskHandlerV2)
	http.Handle(fmt.Sprintf("%s/v2/%s", apiBasePath,  scheduledTasksBasePath), cors.Middleware(handleTaskV2))
	http.Handle(fmt.Sprintf("%s/v2/%s/", apiBasePath,  scheduledTasksBasePath), cors.Middleware(handleTasksV2))
}

func tasksHandlerV2(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		httphelper.GetRequest(w, r, "v2/scheduled-task/:id", func(params map[string]string) (interface{}, error) {
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

func taskHandlerV2(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		httphelper.GetQueryRequest(w, r,  func(values url.Values) (interface{}, error) {
			dayText := values.Get("day")
			day, err := time.Parse("2006-01-02", dayText)
			if err != nil { return nil, err }
			return getSchedulesForDate(day)
		})
	case http.MethodOptions:
		return
	}
}
