package repo_implementation

import (
	"errors"
	"fmt"
	"lab2/src/database/pgdb"
	"lab2/src/reservation/models"
	"lab2/src/utility"
)

const (
	insertReservation      = `insert into reservation(reservation_uid, username, book_uid, library_uid, status, start_date, till_date) values($1, $2, $3, $4, $5, $6, $7) returning id;`
	countUserReservations  = `select count(*) from reservation where username=$1 and status like $2;`
	selectReservation      = `select id, reservation_uid, username, book_uid, library_uid, status, to_char(start_date,'YYYY-MM-DD'), to_char(till_date,'YYYY-MM-DD') from reservation where reservation_uid=$1;`
	selectUserReservations = `select id, reservation_uid, username, book_uid, library_uid, status, to_char(start_date,'YYYY-MM-DD'), to_char(till_date,'YYYY-MM-DD') from reservation where username=$1;`
	updateReservation      = `update reservation set status=$2 where reservation_uid=$1`
)

type ReservationRepository struct {
	db *pgdb.DBManager
}

func NewReservationRepository(manager *pgdb.DBManager) *ReservationRepository {
	return &ReservationRepository{db: manager}
}

func (rr *ReservationRepository) Create(item *models.Reservation) error {
	data, err := rr.db.Query(insertReservation, item.ReservationUid, item.Username, item.BookUid, item.LibraryUid, item.Status, item.StartDate, item.TillDate)
	if err != nil {
		fmt.Printf("Failed to insert a reservation for user %s in db\n", item.Username)
		return err
	}
	if len(data) == 0 {
		fmt.Printf("No id was returned by inserting a reservation for user %s in db\n", item.Username)
		return errors.New("Cannot create person in database")
	}
	return err
}

func (rr *ReservationRepository) GetReservation(ruid string) (*models.Reservation, error) {
	data, err := rr.db.Query(selectReservation, ruid)
	if err != nil {
		fmt.Printf("Failed to get libraries from db\n")
		return nil, err
	}
	res := &models.Reservation{
		Id:             utility.BytesToInt(data[0][0]),
		ReservationUid: utility.BytesToUid(data[0][1]),
		Username:       utility.BytesToString(data[0][2]),
		BookUid:        utility.BytesToUid(data[0][3]),
		LibraryUid:     utility.BytesToUid(data[0][4]),
		Status:         utility.BytesToString(data[0][5]),
		StartDate:      utility.BytesToString(data[0][6]),
		TillDate:       utility.BytesToString(data[0][7]),
	}

	return res, err
}

func (rr *ReservationRepository) GetUserReservationsCount(username string, status string) (int, error) {
	res := -1
	if status == "" {
		status = "%"
	}
	data, err := rr.db.Query(countUserReservations, username, status)
	if err != nil {
		fmt.Printf("Failed to get libraries from db\n")
	} else {
		res = utility.BytesToInt(data[0][0])
	}
	return res, err
}

func (rr *ReservationRepository) GetUserReservations(username string) ([]*models.Reservation, error) {
	data, err := rr.db.Query(selectUserReservations, username)
	if err != nil {
		fmt.Printf("Failed to get libraries from db\n")
		return nil, err
	}
	var res []*models.Reservation

	for _, row := range data {
		res = append(res, &models.Reservation{
			Id:             utility.BytesToInt(row[0]),
			ReservationUid: utility.BytesToUid(row[1]),
			Username:       utility.BytesToString(row[2]),
			BookUid:        utility.BytesToUid(row[3]),
			LibraryUid:     utility.BytesToUid(row[4]),
			Status:         utility.BytesToString(row[5]),
			StartDate:      utility.BytesToString(row[6]),
			TillDate:       utility.BytesToString(row[7]),
		})
	}

	return res, err
}

func (rr *ReservationRepository) UpdateReservation(uid string, status string) error {
	affected, err := rr.db.Exec(updateReservation, uid, status)
	if err != nil {
		fmt.Printf("Failed to update reservation with uid %s in db\n", uid)
		return err
	}
	if affected == 0 {
		err = errors.New("Found no records for the given username")
	}
	return err
}
