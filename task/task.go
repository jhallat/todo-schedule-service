package task

type ScheduledTask struct {
	ScheduleId      int    `json:"scheduleId"`
	TaskId          int    `json:"taskId"`
	TaskDescription string `json:"taskDescription"`
	TaskQuantity	int    `json:"taskQuantity"`
}

type UpdatedTask struct {
	 Id          int    `json:"id"`
	 Description string `json:"description"`
}