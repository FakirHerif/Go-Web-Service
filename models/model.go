package models

import (
	"database/sql"
	"strconv"

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

// @Summary Get a list of 20 persons
// @Description Get persons list from the database
// @Tags persons
// @Accept json
// @Produce json
// @Success 200 {object} Person
// @Router /api/v1/person [get]
func GetPersons(count int) ([]Person, error) {

	rows, err := DB.Query("SELECT id, first_name, last_name, email, ip_address from people LIMIT " + strconv.Itoa(count))

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
// @Tags persons
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
// @Tags persons
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
// @Tags persons
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
// @Tags persons
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
