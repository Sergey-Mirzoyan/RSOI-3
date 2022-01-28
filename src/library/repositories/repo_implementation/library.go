package repo_implementation

import (
	"fmt"
	"lab2/src/database/pgdb"
	"lab2/src/library/models"
	"lab2/src/utility"
)

const (
	selectLibrariesByCity = "select id, library_uid, name, city, address from library where city=$1;"
	selectLibrariesByUid  = "select id, library_uid, name, city, address from library where library_uid=$1;"
)

type LibraryRepository struct {
	db *pgdb.DBManager
}

func NewLibraryRepository(manager *pgdb.DBManager) *LibraryRepository {
	return &LibraryRepository{db: manager}
}

func (lr *LibraryRepository) GetByUid(luid string) (*models.Library, error) {
	data, err := lr.db.Query(selectLibrariesByUid, luid)
	if err != nil {
		fmt.Printf("Failed to get libraries from db\n")
	}

	res := &models.Library{
		Id:         utility.BytesToInt(data[0][0]),
		LibraryUid: utility.BytesToUid(data[0][1]),
		Name:       utility.BytesToString(data[0][2]),
		City:       utility.BytesToString(data[0][3]),
		Address:    utility.BytesToString(data[0][4]),
	}
	return res, err
}

func (lr *LibraryRepository) GetByCity(city string) ([]*models.Library, error) {
	data, err := lr.db.Query(selectLibrariesByCity, city)
	if err != nil {
		fmt.Printf("Failed to get libraries from db\n")
	}
	var res []*models.Library

	for _, row := range data {
		res = append(res, &models.Library{
			Id:         utility.BytesToInt(row[0]),
			LibraryUid: utility.BytesToUid(row[1]),
			Name:       utility.BytesToString(row[2]),
			City:       utility.BytesToString(row[3]),
			Address:    utility.BytesToString(row[4]),
		})
	}
	return res, err
}
