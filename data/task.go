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
func (subTask *SubTask) Create() (err error) {
	statement := "insert into subtasks(uuid, title, completed, task_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)"

	stmt, err := Db.Prepare(statement)

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	stmt.QueryRow(createUUID(), subTask.Title, subTask.Completed, subTask.TaskId, time.Now(), time.Now()).Scan(&subTask.Id, &subTask.Uuid, &subTask.CreatedAt, &subTask.UpdatedAt)

	return
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

func UpdateSubTaskCompletion(completed bool, subTaskId string, taskId string) (err error) {
	statement := "update subtasks set completed = $3 where id = $1 AND task_id = $2"

	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(subTaskId, taskId, completed)

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

func DeleteSubTask(id string, taskId string) (err error) {
	statement := "delete from subtasks WHERE id = $1 AND task_id = $2"

	stmt, err := Db.Prepare(statement)

	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(id, taskId)

	if err != nil {
		return
	}

	return
}

// SELECT * FROM albums JOIN artists ON albums.artist_id = artists.id;

type TaskWithSubTasks struct {
	Id        int
	Uuid      string
	Title     string
	Completed bool
	Deadline  time.Time
	ColorTag  string
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
	SubTasks  []SubTask
}

// get all tasks for a user
func GetTasksByUserId(userId int) (tasks []TaskWithSubTasks, err error) {
	rows, err := Db.Query("SELECT * FROM tasks WHERE user_id = $1", userId)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var task TaskWithSubTasks
		if err = rows.Scan(&task.Id, &task.Uuid, &task.Title, &task.Completed, &task.Deadline, &task.ColorTag, &task.UserId, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return
		}

		rows, err := Db.Query("SELECT * FROM subtasks WHERE task_id = $1", task.Id)

		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		var subTasks []SubTask

		for rows.Next() {
			var subTask SubTask

			rows.Scan(&subTask.Id, &subTask.Uuid, &subTask.Title, &subTask.Completed, &subTask.TaskId, &subTask.CreatedAt, &subTask.UpdatedAt)

			// err != nil {
			// 	// return
			// }

			subTasks = append(subTasks, subTask)
		}

		task.SubTasks = subTasks

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
