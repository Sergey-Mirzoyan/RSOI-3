package repositories

import (
	"encoding/binary"
	"errors"
	"fmt"
	"lab1/database/pgdb"
	"lab1/models"
)

type IPersonsRepository interface {
	Create(person *models.InputPerson) (int, error)
	GetById(id int) (*models.Person, error)
	GetAll() ([]*models.Person, error)
	Update(id int, person *models.InputPerson) (int, error)
	Delete(id int) (int, error)
}

const (
	insertPerson    = "insert into persons(name, age, address, work) values($1, $2, $3, $4) returning id;"
	selectPerson    = "select name, age, address, work from persons where id = $1;"
	selectAllPeople = "select id, name, age, address, work from persons;"
	updatePeople    = "update persons set name=$2, age=$3, address=$4, work=$5 where id = $1;"
	deletePerson    = "delete from persons where id = $1;"
)

type PersonsRepository struct {
	db *pgdb.DBManager
}

func BytesToInt(data []byte) int {
	return int(binary.BigEndian.Uint32(data))
}

func BytesToString(data []byte) string {
	return string(data)
}

func NewPersonsRepository(manager *pgdb.DBManager) *PersonsRepository {
	return &PersonsRepository{db: manager}
}

func (pr *PersonsRepository) Create(person *models.InputPerson) (int, error) {
	data, err := pr.db.Query(insertPerson, person.Name, person.Age, person.Address, person.Work)
	if err != nil {
		fmt.Printf("Failed to insert a person named %s in db\n", *(person.Name))
		return 0, err
	}
	if len(data) == 0 {
		fmt.Printf("No id was returned by inserting a person named %s in db\n", *(person.Name))
		return 0, errors.New("Cannot create person in database")
	}
	return BytesToInt(data[0][0]), err
}

func (pr *PersonsRepository) GetById(id int) (*models.Person, error) {
	data, err := pr.db.Query(selectPerson, id)
	if err != nil {
		fmt.Printf("Failed to get a person with id=%d from db\n", id)
	}
	var res *models.Person
	switch len(data) {
	case 0:
		err = errors.New("Failed to find a person with the given ID")
		res = nil
	case 1:
		res = &models.Person{
			Id:      id,
			Name:    BytesToString(data[0][0]),
			Age:     BytesToInt(data[0][1]),
			Address: BytesToString(data[0][2]),
			Work:    BytesToString(data[0][3]),
		}
	default:
		err = errors.New("Database failure: not unique ID")
		res = nil
	}
	return res, err
}

func (pr *PersonsRepository) GetAll() ([]*models.Person, error) {
	data, err := pr.db.Query(selectAllPeople)
	if err != nil {
		fmt.Printf("Failed to get people from db\n")
	}
	var res []*models.Person

	for _, row := range data {
		res = append(res, &models.Person{
			Id:      BytesToInt(row[0]),
			Name:    BytesToString(row[1]),
			Age:     BytesToInt(row[2]),
			Address: BytesToString(row[3]),
			Work:    BytesToString(row[4]),
		})
	}
	return res, err
}

func (pr *PersonsRepository) Update(id int, person *models.InputPerson) (int, error) {
	affected, err := pr.db.Exec(updatePeople, id, person.Name, person.Age, person.Address, person.Work)
	if err != nil {
		fmt.Printf("Failed to update a person with id=%d in db\n", id)
		return -1, err
	}
	if affected == 0 {
		err = errors.New("No person found in database that matches the received id")
	}
	return affected, err
}

func (pr *PersonsRepository) Delete(id int) (int, error) {
	affected, err := pr.db.Exec(deletePerson, id)
	if err != nil {
		fmt.Printf("Failed to delete a person with id=%d in db\n", id)
		return -1, err
	}
	if affected == 0 {
		err = errors.New("No person found in database that matches the received id")
	}
	return affected, err
}
