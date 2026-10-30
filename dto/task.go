package dto

type TaskRequest struct {
	Subject     string `json:"subject"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Duration    int64  `json:"duration"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
}
