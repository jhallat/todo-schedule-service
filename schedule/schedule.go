package schedule

import (
	"fmt"
	"strings"
	"time"
)

type JsonDate time.Time

type UpdatedTask struct {
	 Id          int    `json:"id"`
	 Description string `json:"description"`
}

type Weekly struct {
	Sunday           int    `json:"sunday"`
	Monday           int    `json:"monday"`
	Tuesday          int    `json:"tuesday"`
	Wednesday        int    `json:"wednesday"`
	Thursday         int    `json:"thursday"`
	Friday           int    `json:"friday"`
	Saturday         int    `json:"saturday"`
}

type Daily struct {
	Day	     JsonDate  `json:"day"`
	Quantity int       `json:"quantity"`
}

type Schedule struct {
	TaskId           int         `json:"taskId"`
	Description      string      `json:"description"`
	Quantifiable     bool        `json:"quantifiable"`
	ScheduleType	 string      `json:"scheduleType"`
	Paused           bool        `json:"paused"`
	GoalId			 int         `json:"goalId"`
	GoalDescription  string      `json:"goalDescription"`
	Weekly  	     *Weekly     `json:"weekly"`
	Daily   	     *Daily      `json:"daily"`
}

type ScheduleForDay struct {
	TaskId           int         `json:"taskId"`
	Description      string      `json:"description"`
	Quantifiable     bool        `json:"quantifiable"`
	ScheduleType	 string      `json:"scheduleType"`
	Quantity	     int         `json:"quantity"`
	GoalId			 int         `json:"goalId"`
	GoalDescription  string      `json:"goalDescription"`
	Paused           bool        `json:"paused"`
}

type SchedulesForRange struct {
	Day	           JsonDate          `json:"day"`
	ScheduledTasks *[]ScheduleForDay `json:"scheduledTasks""`
}

func (j *JsonDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonDate(t)
	return nil
}

func (j JsonDate) MarshalJSON() ([]byte, error) {
	return []byte(j.String()), nil
}

func (j *JsonDate) String() string {
	t := time.Time(*j)
	return fmt.Sprintf("%q", t.Format("2006-01-02"))
}

func (j JsonDate) Date() (year int, month time.Month, day int) {
	t := time.Time(j)
	return t.Date()
}