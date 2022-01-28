package uc_implementation

import (
	"github.com/google/uuid"
	"lab2/src/reservation/models"
	"lab2/src/reservation/repositories"
	"time"
)

type ReservationUsecase struct {
	rr repositories.IReservationRepository
}

func NewReservationUsecase(repo repositories.IReservationRepository) *ReservationUsecase {
	return &ReservationUsecase{rr: repo}
}

func (rc *ReservationUsecase) GetReservation(ruid string) (*models.Reservation, error) {
	return rc.rr.GetReservation(ruid)
}

func (rc *ReservationUsecase) GetUserReservationsCount(username string, status string) (int, error) {
	return rc.rr.GetUserReservationsCount(username, status)
}

func (rc *ReservationUsecase) GetUserReservations(username string) ([]*models.Reservation, error) {
	return rc.rr.GetUserReservations(username)
}

func (rc *ReservationUsecase) Create(item *models.InputReservation) (string, error) {
	uid := uuid.New().String()
	reservation := &models.Reservation{
		Id:             0,
		ReservationUid: uid,
		Username:       item.Username,
		BookUid:        item.BookUid,
		LibraryUid:     item.LibraryUid,
		Status:         "RENTED",
		StartDate:      time.Now().Format("2006-01-02"),
		TillDate:       item.TillDate,
	}
	return uid, rc.rr.Create(reservation)
}

func (rc *ReservationUsecase) UpdateReservation(uid string, status string) error {
	return rc.rr.UpdateReservation(uid, status)
}
