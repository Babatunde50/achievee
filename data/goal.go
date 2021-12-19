package data

import "time"

type Goal struct {
	Id                    int       `json:"id"`
	Uuid                  string    `json:"uuid"`
	Title                 string    `json:"title"`
	ColorTag              string    `json:"colorTag"`
	TotalPercent          int       `json:"totalPercent"`
	TotalPercentCompleted int       `json:"totalPercentCompleted"`
	Completed             bool      `json:"completed"`
	Paused                bool      `json:"paused"`
	Deadline              time.Time `json:"deadline"`
	UserId                int       `json:"userId"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}
