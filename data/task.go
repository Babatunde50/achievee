package data

import (
	"time"
)

type Task struct {
	Id        int       `json:"id"`
	Uuid      string    `json:"uuid"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	Deadline  time.Time `json:"deadline"`
	ColorTag  string    `json:"colorTag"`
	UserId    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SubTask struct {
	Id        int       `json:"id"`
	Uuid      string    `json:"uuid"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	TaskId    int       `json:"taskId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
