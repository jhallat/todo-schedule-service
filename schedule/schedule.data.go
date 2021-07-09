package schedule

import (
	"database/sql"
	"errors"
	"github.com/jhallat/todo-schedule-service/database"
	"log"
	"strings"
)

type WeeklyScheduleDTO struct {
	Id           int          `json:"id"`
	Description  string       `json:"description"`
	Days         string       `json:"selectedDays"`
}

func convertStringToDays(dayString string) (DaySelection, error){
	var converted DaySelection
	var days = strings.Split(dayString, "")
	if len(days) != 7 {
		return converted, errors.New("invalid day string")
	}
	converted.Sunday = days[0] == "1"
	converted.Monday = days[1] == "1"
	converted.Tuesday = days[2] == "1"
	converted.Wednesday = days[3] == "1"
	converted.Thursday = days[4] == "1"
	converted.Friday = days[5] == "1"
	converted.Saturday = days[6] == "1"
	return converted, nil
}

func convertDaysToString(days DaySelection) string {
	converted := ""
	if days.Sunday {
		converted = converted + "1"
	} else {
		converted = converted + "0"
	}
	if days.Monday {
		converted = converted + "1"
	} else {
		converted = converted + "0"
	}
	if days.Tuesday {
		converted = converted + "1"
	} else {
		converted = converted + "0"
	}
	if days.Wednesday {
		converted = converted + "1"
	} else {
		converted = converted + "0"
	}
	if days.Thursday {
		converted = converted + "1"
	} else {
		converted = converted + "0"
	}
	if days.Friday {
		converted = converted + "1"
	} else {
		converted = converted + "0"
	}
	if days.Saturday {
		converted = converted + "1"
	} else {
		converted = converted + "0"
	}
	return converted
}


func getSchedule(id int) (*WeeklySchedule, error) {
	row := database.DbConnection.QueryRow(`SELECT id, 
       description, 
       days
       FROM weeklyschedule
       WHERE id = ?`, id)
	schedule := &WeeklyScheduleDTO{}
	err := row.Scan(&schedule.Id,
		&schedule.Description,
		&schedule.Days)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	selectedDays, _ := convertStringToDays(schedule.Days)
	convertedSchedule := &WeeklySchedule{
		Id: schedule.Id,
		Description: schedule.Description,
		SelectedDays: selectedDays,
	}
	return convertedSchedule, nil
}

func removeSchedule(id int) error {
	_, err := database.DbConnection.Query(`DELETE FROM weeklyschedule 
       WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return nil
}

func getScheduleList() ([]WeeklySchedule, error) {
	results, err := database.DbConnection.Query(`SELECT id, 
       description, 
       days as selectedDays 
       FROM weeklyschedule`)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	schedules := make([]WeeklySchedule, 0)
	for results.Next() {
		var schedule WeeklyScheduleDTO
		results.Scan(&schedule.Id,
			         &schedule.Description,
			         &schedule.Days)
		selectedDays, _ := convertStringToDays(schedule.Days)
		convertedSchedule := WeeklySchedule{
			Id: schedule.Id,
			Description: schedule.Description,
			SelectedDays: selectedDays,
		}
		schedules = append(schedules, convertedSchedule)
	}
	return schedules, nil
}


func updateSchedule(schedule WeeklySchedule) error {
	days := convertDaysToString(schedule.SelectedDays)
	_, err := database.DbConnection.Exec(`UPDATE weeklyschedule 
    SET description = ?,
        days = ?
    WHERE id = ?`,
		schedule.Description,
		days,
		schedule.Id)
	if err != nil {
		return err
	}
	return nil
}

func insertSchedule(schedule WeeklySchedule) (int, error) {
	days := convertDaysToString(schedule.SelectedDays)
	result, err := database.DbConnection.Exec(`INSERT INTO 
        weeklyschedule (description, days)  
        VALUES ($1, $2);`,
		schedule.Description,
		days)
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
