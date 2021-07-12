package task

type ScheduledTask struct {
	ScheduleId      int    `json:"scheduleId"`
	TaskId          int    `json:"taskId""`
	TaskDescription string `json:"taskDescription"`
	TaskQuantity	int    `json:"taskQuantity"`
}
