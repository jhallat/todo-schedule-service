package schedule

import (
	"fmt"
	"github.com/jhallat/todo-schedule-service/cors"
	"github.com/jhallat/todo-schedule-service/httphelper"
	"net/http"
	"strconv"
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
		httphelper.GetRequest(w,  r, "", func(map[string]string) (interface{}, error){ return getScheduleList() })
	case http.MethodPost:
		var weeklySchedule WeeklySchedule
		httphelper.PostRequest(w, r, "Id", &weeklySchedule, func( ) (int64, error) {
			id, err := insertSchedule(weeklySchedule)
			return int64(id), err
		})
	case http.MethodOptions:
		return
	}
}



func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		httphelper.GetRequest(w,  r, "schedules/:id", func(params map[string]string) (interface{}, error){
			id := params["id"]
			scheduleId, err := strconv.Atoi(id)
			if err != nil { return nil, err }
			return getSchedule(scheduleId)
		})
	case http.MethodPut:
		var weeklySchedule WeeklySchedule
		httphelper.PutRequest(w,  r,
			"schedules/:id",
			"Id",
			weeklySchedule,
			func(params map[string]string) (interface{}, error){
				id := params["id"]
				scheduleId, err := strconv.Atoi(id)
				if err != nil { return nil, err }
				return getSchedule(scheduleId)
		     },
		     func() error {
				return updateSchedule(weeklySchedule)
			 })
	case http.MethodOptions:
		return
	case http.MethodDelete:
		httphelper.DeleteRequest(w, r,
			"schedule/:id",
			func(params map[string]string) error {
				id := params["id"]
				scheduleId, err := strconv.Atoi(id)
				if err != nil { return err}
				removeSchedule(scheduleId)
				return nil
			})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
