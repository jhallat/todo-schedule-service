package health

type Check struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Health struct {
	Status string `json:"status"`
	Checks []Check `json:"checks"`
}
