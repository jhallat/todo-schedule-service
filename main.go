package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jhallat/todo-schedule-service/config"
	"github.com/jhallat/todo-schedule-service/database"
	"github.com/jhallat/todo-schedule-service/health"
	"github.com/jhallat/todo-schedule-service/logger"
	"github.com/jhallat/todo-schedule-service/schedule"
	"net/http"
)

type Config struct {
	DbHost     string `prop:"database.host" env:"SCHED_DB_HOST"`
	DbUser     string `prop:"database.user" env:"SCHED_DB_USER"`
	DbPassword string `prop:"database.password" env:"SCHED_DB_PASSWORD"`
	DbPort     string `prop:"database.port" env:"SCHED_DB_PORT"`
	DbName	   string `prop:"database.name" env:"SCHED_DB_NAME"`
	QueueUrl   string `prop:"queue.url" env:"SCHED_QUEUE_URL"`
}

const apiBasePath = "/api"

func main() {

	logger.LogMessage(logger.INFO, "Schedule API starting.")
	configuration := &Config{}
	config.Scan(configuration)
	connection := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		configuration.DbUser, configuration.DbPassword, configuration.DbHost, configuration.DbPort, configuration.DbName)
	database.SetupDatabase(connection)
	schedule.SetupRoutes(apiBasePath)
	health.SetupHealth()
	schedule.SetupTaskListener(configuration.QueueUrl)
	err := http.ListenAndServe(":5002", nil)
	if err != nil {
		fmt.Println(err)
	}
}
