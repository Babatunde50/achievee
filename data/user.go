package data

import (
	"database/sql"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int       `json:"id"`
	Uuid      string    `json:"uuid"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

// create a new user
func (user *User) Create() (err error) {
	statement := "insert into users (uuid, firstName, lastName, email, password, created_at) values ($1, $2, $3, $4, $5, $6)"

	stmt, err := Db.Prepare(statement)

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	stmt.QueryRow(createUUID(), user.FirstName, user.LastName, user.Email, user.Password, time.Now()).Scan(&user.Id, &user.Uuid, &user.CreatedAt)

	// TODO: Manipulate user and add id and uuid field to it..
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

// Check if email already exists
func UserEmailExists(email string) (bool, error) {

	user := User{}

	err := Db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&user.Email)

	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return true, err
	}

	return true, nil
}

// Get a single user given the UUID
func UserByUUID(uuid string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).
		Scan(&user.Id, &user.Uuid, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}

//Generate a salted hash for the input string
func GeneratePassword(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

//Compare string to generated hash
func ComparePasswords(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
