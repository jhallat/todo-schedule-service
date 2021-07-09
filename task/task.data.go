package task

import "github.com/jhallat/todo-schedule-service/database"

func insertScheduledTask(task ScheduledTask) error {
	_, err := database.DbConnection.Exec(`INSERT INTO
     scheduledtask (schedule_id, task_id, task_description)
     values ($1, $2, $3)`,
     task.ScheduleId,
     task.TaskId,
     task.TaskDescription)
	return err
}

func getScheduledTasksBySchedule(scheduleId int) ([]ScheduledTask, error) {
	results, err := database.DbConnection.Query(`SELECT schedule_id as scheduleId,
            task_id as taskId,
            task_description as TaskDescription`)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	tasks := make([]ScheduledTask, 0)
	for results.Next() {
		var task ScheduledTask
		results.Scan(&task.ScheduleId,
			         &task.TaskId,
			         &task.TaskDescription)
		tasks = append(tasks, task)
	}
	return tasks, nil
}
