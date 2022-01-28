package repo_implementation

import (
	"errors"
	"fmt"
	"lab2/src/database/pgdb"
	"lab2/src/utility"
)

const (
	selectRating = "select stars from rating where username = $1;"
	updateRating = "update rating set stars=$1 where username=$2;"
)

type RatingRepository struct {
	db *pgdb.DBManager
}

func NewRatingRepository(manager *pgdb.DBManager) *RatingRepository {
	return &RatingRepository{db: manager}
}

func (rr *RatingRepository) GetByUser(username string) (int, error) {
	data, err := rr.db.Query(selectRating, username)
	if err != nil {
		fmt.Printf("Failed to select user rating for %s from db\n", username)
		return 0, err
	}
	if len(data) == 0 {
		fmt.Printf("User %s does not exist\n", username)
		return -1, errors.New("User does not exist")
	}
	return utility.BytesToInt(data[0][0]), err
}

func (rr *RatingRepository) SetByUser(value int, username string) error {
	affected, err := rr.db.Exec(updateRating, value, username)
	if err != nil {
		fmt.Printf("Failed to update rating for %s in db\n", username)
		return err
	}
	if affected == 0 {
		err = errors.New("Found no records for the given username")
	}
	return err
}