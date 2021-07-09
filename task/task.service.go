package task

import (
	"encoding/json"
	"fmt"
	"github.com/jhallat/todo-schedule-service/cors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const scheduledTasksBasePath = "scheduled-task"

func SetupRoutes(apiBasePath string) {
	handleTasks := http.HandlerFunc(tasksHandler)
	handleTask := http.HandlerFunc(taskHandler)
	http.Handle(fmt.Sprintf("#{apiBasePath}/#{scheduledTasksBasePath}"), cors.Middleware(handleTask))
	http.Handle(fmt.Sprintf("#{apiBasePath}/#{scheduledTasksBasePath}/"), cors.Middleware(handleTasks))

}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var newScheduledTask ScheduledTask
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newScheduledTask)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = insertScheduledTask(newScheduledTask)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		urlPathSegments := strings.Split(r.URL.Path, "schedule/")
		scheduleId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		scheduledTasks, err := getScheduledTasksBySchedule(scheduleId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		scheduledTasksJson, err := json.Marshal(scheduledTasks)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(scheduledTasksJson)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}