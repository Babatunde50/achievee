package data

import (
	"fmt"
	"log"
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

// create a new task
func (task *Task) Create() (err error) {
	statement := "insert into tasks(uuid, title, completed, deadline, color_tag, user_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8)"

	stmt, err := Db.Prepare(statement)

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	fmt.Println(task.Deadline, "task.Deadline!!")

	stmt.QueryRow(createUUID(), task.Title, task.Completed, task.Deadline, task.ColorTag, task.UserId, time.Now(), time.Now()).Scan(&task.Id, &task.Uuid, &task.CreatedAt, &task.UpdatedAt)

	return
}

func (task *Task) Update() (err error) {
	statement := "update tasks set title = $3, completed = $4, deadline = $5, color_tag = $6 where id = $1 AND user_id = $2"

	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Id, task.UserId, task.Title, task.Completed, task.Deadline, task.ColorTag)
	return
}

type Completed struct {
	IsCompleted bool
	UserId      int
	TaskId      int
}

func UpdateTaskCompletion(completed Completed) (err error) {
	statement := "update tasks set completed = $3 where id = $1 AND user_id = $2"

	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(completed.TaskId, completed.UserId, completed.IsCompleted)

	return
}

type UpdateTaskData struct {
	UserId   int
	TaskId   int
	Title    string
	ColorTag string
	Deadline time.Time
}

func UpdateTaskEdits(editData UpdateTaskData) (err error) {
	statement := "update tasks set title = $3, deadline = $4, color_tag = $5 where id = $1 AND user_id = $2"

	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(editData.TaskId, editData.UserId, editData.Title, editData.Deadline, editData.ColorTag)

	return
}

// delete a task
func (task *Task) Delete() (err error) {
	statement := "delete from tasks WHERE id = $1 AND user_id = $2"

	stmt, err := Db.Prepare(statement)

	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(task.Id, task.UserId)

	if err != nil {
		return
	}

	return
}

// get all tasks for a user
func GetTasksByUserId(userId int) (tasks []Task, err error) {
	rows, err := Db.Query("SELECT * FROM tasks WHERE user_id = $1", userId)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rows, "this is the rows from query")

	defer rows.Close()

	for rows.Next() {
		var task Task
		if err = rows.Scan(&task.Id, &task.Uuid, &task.Title, &task.Completed, &task.Deadline, &task.ColorTag, &task.UserId, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
