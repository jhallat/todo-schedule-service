package task

import (
	"github.com/jhallat/todo-schedule-service/database"
)

const sqlFindBySchedule = `
    	SELECT schedule_id as scheduleId,
               task_id as taskId,
               task_description as TaskDescription,
               task_quantity as TaskQuantity
        FROM scheduled_task
        WHERE schedule_id = $1`

const sqlInsertTask = `
     INSERT INTO
     scheduled_task (schedule_id, task_id, task_description, task_quantity)
     values ($1, $2, $3, $4)`

func insertScheduledTask(task ScheduledTask) error {
	_, err := database.DbConnection.Exec(sqlInsertTask,
     task.ScheduleId,
     task.TaskId,
     task.TaskDescription,
     task.TaskQuantity)
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
			         &task.TaskQuantity)
		tasks = append(tasks, task)
	}
	return tasks, nil
}
