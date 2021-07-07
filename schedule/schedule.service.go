package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/jhallat/todo-schedule-service/cors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const schedulesBasePath = "schedules"

func SetupRoutes(apiBasePath string) {
	handleSchedules := http.HandlerFunc(schedulesHandler)
	handleSchedule := http.HandlerFunc(scheduleHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, schedulesBasePath), cors.Middleware(handleSchedules))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, schedulesBasePath), cors.Middleware(handleSchedule))
}

func schedulesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		weeklyScheduleList, err := getScheduleList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		weeklyScheduleJson, err := json.Marshal(weeklyScheduleList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(weeklyScheduleJson)
	case http.MethodPost:
		var newWeeklySchedule WeeklySchedule
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newWeeklySchedule)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newWeeklySchedule.Id != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = addOrUpdateSchedule(newWeeklySchedule)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	}
}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "schedule/")
	scheduleId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	schedule := getSchedule(scheduleId)
	if schedule == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		scheduleJSON, err := json.Marshal(schedule)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(scheduleJSON)
	case http.MethodPut:
		var updatedSchedule WeeklySchedule
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedSchedule)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedSchedule.Id != scheduleId {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		addOrUpdateSchedule(updatedSchedule)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodOptions:
		return
	case http.MethodDelete:
		removeSchedule(scheduleId)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
