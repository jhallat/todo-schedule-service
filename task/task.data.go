package task

import (
	"database/sql"
	"errors"
	"github.com/jhallat/todo-schedule-service/database"
	"time"
)

const sqlFindBySchedule = `
    	SELECT schedule_id as scheduleId,
               task_id as taskId,
               task_description as TaskDescription,
               task_quantity as TaskQuantity,
    	       goal_id as GoalId,
    	       goal_description as GoalDescription
        FROM scheduled_task
        WHERE schedule_id = $1`

const sqlInsertTask = `
     INSERT INTO
     scheduled_task (schedule_id, task_id, task_description, task_quantity, goal_id, goal_description)
     values ($1, $2, $3, $4, $5, $6)`

const sqlUpdateTask = `
	UPDATE scheduled_task
    SET task_description = $2
    WHERE task_id = $1
`


const sqlInsertSchedule = `
	INSERT INTO
	schedule (task_id, description, quantifiable, schedule_type, paused, goal_id, goal_description)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
`

const sqlUpdateSchedule = `
	UPDATE schedule
    SET schedule_type = $2,
        paused = $3,
        description = $4,
        quantifiable = $5,
        goal_id = $6,
        goal_description = $7
    WHERE task_id = $1
`

const sqlGetScheduleCount = `
	SELECT COUNT(*)
    FROM schedule
    WHERE task_id = $1
`

const sqlGetSchedule = `
	SELECT task_id as TaskId,
           description as Description,
           quantifiable as Quantifiable,
           schedule_type as ScheduleType,
           paused as Paused,
	       goal_id as GoalId,
	       goal_description as GoalDescription
    FROM schedule
    WHERE task_id = $1
`

const sqlGetAllSchedules = `
	SELECT task_id as TaskId,
           description as Description,
           quantifiable as Quantifiable,
           schedule_type as ScheduleType,
           paused as Paused,
	       goal_id as GoalId,
	       goal_description as GoalDescription
    FROM schedule
    WHERE paused = false
`

const sqlGetWeeklyCount = `
	SELECT COUNT(*)
    FROM weekly
    WHERE task_id = $1
`

const sqlGetWeekly = `
	SELECT sunday as Sunday,
           monday as Monday,
           tuesday as Tuesday,
           wednesday as Wednesday,
           thursday as Thursday,
           friday as Friday,
           saturday as Saturday
    FROM weekly
    WHERE task_id = $1
`

const sqlGetDaily = `
	SELECT day as Day,
           quantity as Quantity
    FROM daily
    WHERE task_id = $1
`


const sqlInsertWeekly = `
	INSERT INTO 
    weekly (task_id, sunday, monday, tuesday, wednesday, thursday, friday, saturday)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

const sqlUpdateWeekly = `
	UPDATE weekly
    SET sunday = $2,
        monday = $3,
        tuesday = $4,
        wednesday = $5,
        thursday = $6,
        friday = $7,
        saturday = $8
    WHERE task_id = $1
`

const sqlDeleteWeekly = `
	DELETE FROM weekly
    WHERE task_id = $1
`

const sqlGetDailyCount = `
	SELECT COUNT(*)
    FROM daily
    WHERE task_id = $1
`

const sqlInsertDaily = `
	INSERT INTO
    daily (task_id, day, quantity)
    values ($1, $2, $3)
`

const sqlUpdateDaily = `
	UPDATE daily
    SET day = $2,
        quantity = $3
    WHERE task_id = $1
`

const sqlDeleteDaily = `
	DELETE FROM daily
    WHERE task_id = $1
`

func getSchedule(taskId int) (*Schedule, error) {
	row := database.DbConnection.QueryRow(sqlGetSchedule, taskId)
	schedule := &Schedule{}
	err := row.Scan(&schedule.TaskId,
		            &schedule.Description,
		            &schedule.Quantifiable,
		            &schedule.ScheduleType,
		            &schedule.Paused,
		            &schedule.GoalId,
		            &schedule.GoalDescription)
	if err == sql.ErrNoRows {
		schedule.TaskId = taskId
		schedule.ScheduleType = "NONE"
		return schedule, nil
	} else if err != nil {
		return nil, err
	}
	if schedule.ScheduleType == "WEEKLY" {
		weekly, err := getWeekly(taskId)
		if err != nil {
			return nil, err
		}
		schedule.Weekly = weekly
		return schedule, nil
	}
	if schedule.ScheduleType == "DAILY" {
		daily, err := getDaily(taskId)
		if err != nil {
			return nil, err
		}
		schedule.Daily = daily
		return schedule, nil
	}
	return schedule, nil
}

func getSchedulesForDate(day time.Time) ([]ScheduleForDay, error) {
	dayOfWeek := day.Weekday()
	results, err := database.DbConnection.Query(sqlGetAllSchedules)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	schedules := make([]ScheduleForDay, 0)
	for results.Next() {
		var schedule Schedule
		results.Scan(&schedule.TaskId,
			         &schedule.Description,
			         &schedule.Quantifiable,
			         &schedule.ScheduleType,
			         &schedule.Paused,
			         &schedule.GoalId,
			         &schedule.GoalDescription)
		scheduleForDay := ScheduleForDay{
			TaskId: schedule.TaskId,
			Description: schedule.Description,
			Quantifiable: schedule.Quantifiable,
			ScheduleType: schedule.ScheduleType,
			GoalId: schedule.GoalId,
			GoalDescription: schedule.GoalDescription,
		}
		if schedule.ScheduleType == "WEEKLY" {
			weekly, err := getWeekly(schedule.TaskId)
			if err != nil {
				return nil, err
			}
			switch dayOfWeek {
			case time.Sunday:
				scheduleForDay.Quantity = weekly.Sunday
			case time.Monday:
				scheduleForDay.Quantity = weekly.Monday
			case time.Tuesday:
				scheduleForDay.Quantity = weekly.Tuesday
			case time.Wednesday:
				scheduleForDay.Quantity = weekly.Wednesday
			case time.Thursday:
				scheduleForDay.Quantity = weekly.Thursday
			case time.Friday:
				scheduleForDay.Quantity = weekly.Friday
			case time.Saturday:
				scheduleForDay.Quantity = weekly.Saturday
			}
		}
		if schedule.ScheduleType == "DAILY" {
			daily, err := getDaily(schedule.TaskId)
			if err != nil {
				return nil, err
			}
			y1, m1, d1 := daily.Day.Date()
			y2, m2, d2 := day.Date()
			if y1 == y2 && m1 == m2 && d1 == d2 {
				scheduleForDay.Quantity = daily.Quantity
			}
		}
		if scheduleForDay.Quantity > 0 {
			schedules = append(schedules, scheduleForDay)
		}
	}
	return schedules, nil
}

func getWeekly(taskId int) (*Weekly, error) {
	row := database.DbConnection.QueryRow(sqlGetWeekly, taskId)
	weekly := &Weekly{}
	err := row.Scan(&weekly.Sunday,
		            &weekly.Monday,
		            &weekly.Tuesday,
		            &weekly.Wednesday,
		            &weekly.Thursday,
		            &weekly.Friday,
		            &weekly.Saturday)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return weekly, nil
}

func getDaily(taskId int) (*Daily, error) {
	row := database.DbConnection.QueryRow(sqlGetDaily, taskId)
	daily := &Daily{}
	err := row.Scan(&daily.Day,
		            &daily.Quantity)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return daily, nil
}

func saveSchedule(schedule Schedule) error {
	row := database.DbConnection.QueryRow(sqlGetScheduleCount, schedule.TaskId)
	var count int
	err := row.Scan(&count)
	if err == nil {
		if count == 0 {
			return insertSchedule(schedule)
		} else {
			return updateSchedule(schedule)
		}
	}
	return err
}

func updateSchedule(schedule Schedule) error {
	_, err := database.DbConnection.Exec(sqlUpdateSchedule,
		schedule.TaskId,
		schedule.ScheduleType,
		schedule.Paused,
		schedule.Description,
		schedule.Quantifiable,
		schedule.GoalId,
		schedule.GoalDescription)
	if err != nil {
		return err
	}
	if schedule.ScheduleType == "WEEKLY" {
		return saveWeekly(schedule.TaskId, schedule.Weekly)
	}
	if schedule.ScheduleType == "DAILY" {
		return saveDaily(schedule.TaskId, schedule.Daily)
	}
	return nil
}

func insertSchedule(schedule Schedule) error {
	_, err := database.DbConnection.Exec(sqlInsertSchedule,
		schedule.TaskId,
		schedule.Description,
		schedule.Quantifiable,
		schedule.ScheduleType,
		schedule.Paused,
		schedule.GoalId,
		schedule.GoalDescription)
	if err != nil {
		return err
	}
	if schedule.ScheduleType == "WEEKLY" {
		return saveWeekly(schedule.TaskId, schedule.Weekly)
	}
	if schedule.ScheduleType == "DAILY" {
		return saveDaily(schedule.TaskId, schedule.Daily)
	}
	return nil
}

func saveWeekly(taskId int, weekly *Weekly) error {
	row := database.DbConnection.QueryRow(sqlGetWeeklyCount, taskId)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = database.DbConnection.Exec(sqlInsertWeekly,
			taskId,
			weekly.Sunday,
			weekly.Monday,
			weekly.Tuesday,
			weekly.Wednesday,
			weekly.Thursday,
			weekly.Friday,
			weekly.Saturday)
		return err
	} else {
		_, err = database.DbConnection.Exec(sqlUpdateWeekly,
			taskId,
			weekly.Sunday,
			weekly.Monday,
			weekly.Tuesday,
			weekly.Wednesday,
			weekly.Thursday,
			weekly.Friday,
			weekly.Saturday)
		return err
	}
}

func saveDaily(taskId int, daily *Daily) error {
	if daily == nil {
		return errors.New("daily value is nil or invalid")
	}
	row := database.DbConnection.QueryRow(sqlGetDailyCount, taskId)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = database.DbConnection.Exec(sqlInsertDaily,
			taskId,
			daily.Day,
			daily.Quantity)
		return err
	} else {
		_, err = database.DbConnection.Exec(sqlUpdateDaily,
			taskId,
			daily.Day,
			daily.Quantity)
		return err
	}
}

func insertScheduledTask(task ScheduledTask) error {
	_, err := database.DbConnection.Exec(sqlInsertTask,
     task.ScheduleId,
     task.TaskId,
     task.TaskDescription,
     task.TaskQuantity,
     task.GoalId,
     task.GoalDescription)
	return err
}

func getScheduledTasksBySchedule(scheduleId int) ([]ScheduledTask, error) {
	results, err := database.DbConnection.Query(sqlFindBySchedule, scheduleId)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	tasks := make([]ScheduledTask, 0)
	for results.Next() {
		var task ScheduledTask
		results.Scan(&task.ScheduleId,
			         &task.TaskId,
			         &task.TaskDescription,
			         &task.TaskQuantity,
			         &task.GoalId,
			         &task.GoalDescription)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func updateDescription(task UpdatedTask) error {
	_, err := database.DbConnection.Exec(sqlUpdateTask,
		task.Id,
		task.Description)
	return err
}
