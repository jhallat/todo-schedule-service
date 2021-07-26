package main

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jhallat/todo-schedule-service/database"
	"github.com/jhallat/todo-schedule-service/health"
	"github.com/jhallat/todo-schedule-service/schedule"
	"github.com/jhallat/todo-schedule-service/task"
	"github.com/jhallat/todo-schedule-service/weeklytask"
	"net/http"
)


const apiBasePath = "/api"
func main() {
	database.SetupDatabase()
	schedule.SetupRoutes(apiBasePath)
	task.SetupRoutes(apiBasePath)
	weeklytask.SetupRoutes(apiBasePath)
	health.SetupHealth()
	task.SetupListener()
	http.ListenAndServe(":5002", nil)
}
