package weeklytask

import (
	"fmt"
	"github.com/jhallat/todo-schedule-service/cors"
	"github.com/jhallat/todo-schedule-service/httphelper"
	"net/http"
	"strconv"
)

const weeklyTasksBasePath = "weekly-tasks"

func SetupRoutes(apiBasePath string) {
	handleWeeklyTasks := http.HandlerFunc(weeklyTasksHandler)
	handleWeeklyTask := http.HandlerFunc(weeklyTaskHandler)
	http.Handle(fmt.Sprintf("%s/%s",apiBasePath,weeklyTasksBasePath), cors.Middleware(handleWeeklyTasks))
	http.Handle(fmt.Sprintf("%s/%s/",apiBasePath,weeklyTasksBasePath), cors.Middleware(handleWeeklyTask))

}

func weeklyTasksHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("weeklyTasksHandler")
	switch r.Method {
	case http.MethodPost:
		fmt.Println("POST")
		var weeklyTask WeeklyTask
		httphelper.PostRequest(w, r, "Id", &weeklyTask, func( ) (int64, error) {
			id, err := insertWeeklyTask(weeklyTask)
			return int64(id), err
		})
	case http.MethodPut:
		fmt.Println("PUT")
		var weeklyTask WeeklyTask
		httphelper.PutRequest(w, r,
			"weekly-tasks/:id",
			"Id",
			weeklyTask,
			func(params map[string]string) (interface{}, error) {
				id := params["id"]
				weeklyTaskId, err := strconv.Atoi(id)
				if err != nil { return nil, err }
				return getWeeklyTask(weeklyTaskId)
			},
			func() error {
				return updateWeeklyTask(weeklyTask)
			})
	case http.MethodOptions:
		return
	}
}

func weeklyTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("weeklyTaskHandler")
	switch r.Method {
	case http.MethodGet:
		httphelper.GetRequest(w, r, "weekly-tasks/:id", func(params map[string]string) (interface{}, error) {
			id := params["id"]
			taskId, err := strconv.Atoi(id)
			if err != nil {
				return nil, err
			}
			return getWeeklyTask(taskId)
		})
	case http.MethodOptions:
		return
	}
}