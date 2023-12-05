package models

import (
	"database/sql"
	"fmt"

	"errors"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

type Person struct {
	Id        int    `json:"id" swaggerignore:"true"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IpAddress string `json:"ip_address"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// @Summary Get a list of 20 persons
// @Description Get persons list from the database
// @Tags person
// @Accept json
// @Produce json
// @Success 200 {object} Person
// @Router /api/v1/person [get]
func GetPersons(limit, offset int) ([]Person, error) {

	query := fmt.Sprintf("SELECT id, first_name, last_name, email, ip_address FROM people LIMIT %d OFFSET %d", limit, offset)

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	people := make([]Person, 0)

	for rows.Next() {
		singlePerson := Person{}
		err = rows.Scan(&singlePerson.Id, &singlePerson.FirstName, &singlePerson.LastName, &singlePerson.Email, &singlePerson.IpAddress)

		if err != nil {
			return nil, err
		}

		people = append(people, singlePerson)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return people, err
}

// @Summary Get a person by ID
// @Description Get a person by their ID from the database
// @Tags person
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} Person
// @Router /api/v1/person/{id} [get]
func GetPersonById(id string) (Person, error) {
	stmt, err := DB.Prepare("SELECT id, first_name, last_name, email, ip_address FROM people WHERE id = ?")

	if err != nil {
		return Person{}, err
	}

	person := Person{}

	sqlErr := stmt.QueryRow(id).Scan(&person.Id, &person.FirstName, &person.LastName, &person.Email, &person.IpAddress)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Person{}, nil
		}
		return Person{}, sqlErr
	}
	return person, nil
}

// @Summary Add a new person
// @Description Add a new person to the database
// @Tags person
// @Accept json
// @Produce json
// @Param person body Person true "New Person Object"
// @Success 200 {string} string "Person added successfully"
// @Router /api/v1/person [post]
func AddPerson(newPerson Person) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO people (first_name, last_name, email, ip_address) VALUES (?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.IpAddress)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

// @Summary Update a person's information by their ID
// @Description Update a person's information in the database by their ID
// @Tags person
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body Person true "Updated Person Object"
// @Success 200 {string} string "Person updated successfully"
// @Router /api/v1/person/{id} [put]
func UpdatePerson(ourPerson Person, id int) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM people WHERE id = ?", id).Scan(&count)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	if count == 0 {
		tx.Rollback()
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE people SET first_name = ?, last_name = ?, email = ?, ip_address = ? WHERE Id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ourPerson.FirstName, ourPerson.LastName, ourPerson.Email, ourPerson.IpAddress, ourPerson.Id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

// @Summary Delete a person by their ID
// @Description Delete a person from the database by their ID
// @Tags person
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {string} string "Person deleted successfully"
// @Router /api/v1/person/{id} [delete]
func DeletePerson(personId int) (bool, error) {
	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM people WHERE id = ?", personId).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		tx.Rollback()
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from people WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(personId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

// @Summary Get a list of users
// @Description Get users list from the database
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Router /api/v1/user [get]
func GetUsers() ([]User, error) {
	rows, err := DB.Query("SELECT id, username, email, password, role FROM user")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return users, nil
}

// @Summary Get a user by ID
// @Description Get a user by their ID from the database
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Router /api/v1/user/{id} [get]
func GetUserByID(userID int) (User, error) {
	var user User
	err := DB.QueryRow("SELECT id, username, email, password, role FROM user WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// @Summary Create a new user
// @Description Create a new user in the database
// @Tags user
// @Accept json
// @Produce json
// @Param newUser body User true "New user details"
// @Success 200 {integer} integer
// @Router /api/v1/user [post]
func CreateUser(newUser User) (int64, error) {
	if newUser.Role != "user" {
		newUser.Role = "user"
	}

	result, err := DB.Exec("INSERT INTO user (username, email, password, role) VALUES (?, ?, ?, ?)", newUser.Username, newUser.Email, newUser.Password, newUser.Role)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// @Summary Update an existing user
// @Description Update an existing user in the database
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID to update"
// @Param updatedUser body User true "Updated user details"
// @Success 200 {string} string
// @Router /api/v1/user/{id} [put]
func UpdateUser(updatedUser User) error {
	if updatedUser.Role == "" {
		updatedUser.Role = "user"
	}

	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM user WHERE id = ?", updatedUser.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("kullanici bulunamadi")
	}

	query := "UPDATE user SET username = ?, email = ?, role = ?"
	var args []interface{}
	args = append(args, updatedUser.Username, updatedUser.Email, updatedUser.Role)

	if updatedUser.Password != "" {
		query += ", password = ?"
		args = append(args, updatedUser.Password)
	}

	query += " WHERE id = ?"
	args = append(args, updatedUser.ID)

	_, err = DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// @Summary Delete a user by ID
// @Description Delete a user from the database by their ID
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID to delete"
// @Success 200 {string} string
// @Router /api/v1/user/{id} [delete]
func DeleteUser(userID int) error {
	result, err := DB.Exec("DELETE FROM user WHERE id = ?", userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("kullanici bulunamadi")
	}

	return nil
}

func GetTotalPersonsCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM people"

	err := DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
