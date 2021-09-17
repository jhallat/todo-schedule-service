package weeklytask

type WeeklyTask struct {
	Id               int    `json:"id"`
	TaskId           int    `json:"taskId"`
	TaskDescription  string `json:"taskDescription"`
	TaskQuantifiable bool   `json:"taskQuantifiable"`
	Paused           bool   `json:"paused"`
	Sunday           int    `json:"sunday"`
	Monday           int    `json:"monday"`
	Tuesday          int    `json:"tuesday"`
	Wednesday        int    `json:"wednesday"`
	Thursday         int    `json:"thursday"`
	Friday           int    `json:"friday"`
	Saturday         int    `json:"saturday"`
}
