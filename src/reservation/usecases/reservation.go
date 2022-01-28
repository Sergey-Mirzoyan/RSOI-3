package usecases

import (
	"lab2/src/reservation/models"
)

type IReservationUsecase interface {
	GetReservation(ruid string) (*models.Reservation, error)
	GetUserReservationsCount(username string, status string) (int, error)
	GetUserReservations(username string) ([]*models.Reservation, error)
	Create(item *models.InputReservation) (string, error)
	UpdateReservation(uid string, status string) error
}
