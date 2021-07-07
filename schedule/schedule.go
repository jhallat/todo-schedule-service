package schedule

type WeeklySchedule struct {
	Id           int          `json:"id"`
	Description  string       `json:"description"`
	SelectedDays DaySelection `json:"selectedDays"`
}

type DaySelection struct {
	Sunday    bool `json:"sunday"`
	Monday    bool `json:"monday"`
	Tuesday   bool `json:"tuesday"`
	Wednesday bool `json:"wednesday"`
	Thursday  bool `json:"thursday"`
	Friday    bool `json:"friday"`
	Saturday  bool `json:"saturday"`
}