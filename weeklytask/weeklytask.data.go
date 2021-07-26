package weeklytask

import (
	"database/sql"
	"github.com/jhallat/todo-schedule-service/database"
	"log"
)

const sqlInsertWeeklyTask = `
	INSERT INTO
	weekly_task (task_id, 
                 task_description, 
                 task_quantifiable,
                 paused,
                 sunday,
                 monday,
                 tuesday,
                 wednesday,
                 thursday,
                 friday,
                 saturday)
    values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

const sqlGetWeeklyTask = `
	SELECT 	id as Id,
	        task_id as TaskId,
	        task_description as TaskDescription,
	        task_quantifiable as TaskQuantifiable,
	        paused as Paused,
	        sunday as Sunday,
	        monday as Monday,
	        tuesday as Tuesday,
	        wednesday as Wednesday,
	        thursday as Thursday,
	        friday as Friday,
	        saturday as Saturday
    FROM weekly_task
    WHERE task_id = $1
`

const sqlUpdateWeeklyTask = `
	UPDATE weekly_task
    SET paused = $2,
        sunday = $3,
        monday = $4,
        tuesday = $5,
        wednesday = $6,
        thursday = $7,
        friday = $8,
        saturday = $9
    WHERE task_id = $1
`

func updateWeeklyTask(weeklyTask WeeklyTask) error {
	_, err := database.DbConnection.Exec(sqlUpdateWeeklyTask,
		weeklyTask.TaskId,
		weeklyTask.Paused,
		weeklyTask.Sunday,
		weeklyTask.Monday,
		weeklyTask.Tuesday,
		weeklyTask.Wednesday,
		weeklyTask.Thursday,
		weeklyTask.Friday,
		weeklyTask.Saturday)
	if err !=nil {
		return err
	}
	return nil
}

func insertWeeklyTask(weeklyTask WeeklyTask) (int, error) {
	result, err := database.DbConnection.Exec(sqlInsertWeeklyTask,
		weeklyTask.Id,
		weeklyTask.TaskDescription,
		weeklyTask.TaskQuantifiable,
		weeklyTask.Paused,
		weeklyTask.Sunday,
		weeklyTask.Monday,
		weeklyTask.Tuesday,
		weeklyTask.Wednesday,
		weeklyTask.Thursday,
		weeklyTask.Friday,
		weeklyTask.Saturday)
	if err != nil {
		log.Print(err)
		return 0, nil
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(insertId), nil
}

func getWeeklyTask(taskId int) (*WeeklyTask, error) {
	row := database.DbConnection.QueryRow(sqlGetWeeklyTask, taskId)
	weeklyTask := WeeklyTask{}
	err := row.Scan(&weeklyTask.Id,
		            &weeklyTask.TaskId,
		            &weeklyTask.TaskDescription,
		            &weeklyTask.TaskQuantifiable,
		            &weeklyTask.Paused,
		            &weeklyTask.Sunday,
		            &weeklyTask.Monday,
		            &weeklyTask.Tuesday,
		            &weeklyTask.Thursday,
		            &weeklyTask.Friday,
		            &weeklyTask.Saturday)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &weeklyTask, nil
}
