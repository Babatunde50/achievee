package data

import (
	"log"
	"time"
)

type User struct {
	Id        int
	Uuid      string
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

// USER DB OPERATIONS

// USERS
// create a new user
func (user *User) Create() (err error) {
	statement := "insert into users (uuid, firstName, lastName, email, password, created_at) values ($1, $2, $3, $4, $5 $6)"

	stmt, err := Db.Prepare(statement)

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	stmt.QueryRow(createUUID(), user.FirstName, user.LastName, user.Email, user.Password, time.Now()).Scan(&user.Id, &user.Uuid, &user.CreatedAt)

	return
}

// Delete user from database
func (user *User) Delete() (err error) {
	statement := "delete from users where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	return
}

// Update user information in the database
func (user *User) Update() (err error) {
	statement := "update users set firstName = $2, lastName= $3 email = $4 where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.FirstName, user.LastName, user.Email)
	return
}

// Get a single user given the email
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, firstName, lastName, email, password, created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}

// Get a single user given the UUID
func UserByUUID(uuid string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).
		Scan(&user.Id, &user.Uuid, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}
