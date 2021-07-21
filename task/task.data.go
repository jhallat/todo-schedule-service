package task

import (
	"github.com/jhallat/todo-schedule-service/database"
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
