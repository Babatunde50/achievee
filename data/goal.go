package data

import (
	"log"
	"time"
)

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

// create goal
func (goal *Goal) Create() (err error) {
	statement := `
		INSERT INTO goals (uuid, title, color_tag, total_percent, total_percent_completed, completed, paused, deadline, user_id, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, uuid, created_at, updated_at
	`

	stmt, err := Db.Prepare(statement)

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	stmt.QueryRow(createUUID(), goal.Title, goal.ColorTag, goal.TotalPercent, goal.TotalPercentCompleted, goal.Completed, goal.Paused, goal.Deadline, goal.UserId, goal.CreatedAt, goal.UpdatedAt).Scan(&goal.Id, &goal.Uuid, &goal.CreatedAt, &goal.UpdatedAt)

	return

}

// delete goal
func (goal *Goal) Delete() (err error) {
	statement := "delete from goals WHERE id = $1 AND user_id = $2"

	stmt, err := Db.Prepare(statement)

	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(goal.Id, goal.UserId)

	return
}

type UpdateGoalData struct {
	UserId       int
	GoalId       int
	Title        string
	ColorTag     string
	TotalPercent int
	Deadline     time.Time
}

// update goal
func UpdateGoal(updateData UpdateGoalData) (err error) {
	statement := "update goals set title = $3, deadline = $4, color_tag = $5, total_percent = $6 where id = $1 AND user_id = $2"

	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updateData.GoalId, updateData.UserId, updateData.Title, updateData.Deadline, updateData.ColorTag, updateData.TotalPercent)

	return
}

type UpdateGoalPercentCompletedData struct {
	UserId                int
	GoalId                int
	TotalPercentCompleted int
}

// update goal progress
func UpdateGoalProgress(updateGoalPercentCompleted UpdateGoalPercentCompletedData) (err error) {

	statement := "update goals set total_percent_completed = $3 where id = $1 AND user_id = $2"

	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updateGoalPercentCompleted.GoalId, updateGoalPercentCompleted.UserId, updateGoalPercentCompleted.TotalPercentCompleted)

	return

}

// get all goals for a user
func GetGoalsByUserId(userId int) (goals []Goal, err error) {
	rows, err := Db.Query("SELECT * FROM tasks WHERE user_id = $1", userId)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var goal Goal

		rows.Scan(&goal.Id, &goal.Uuid, &goal.Title, &goal.ColorTag, &goal.TotalPercent, &goal.TotalPercentCompleted, &goal.Completed, &goal.Paused, &goal.Deadline, &goal.UserId, &goal.CreatedAt, &goal.UpdatedAt)

		goals = append(goals, goal)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
