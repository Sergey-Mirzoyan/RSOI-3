package repositories

import (
	"lab2/src/reservation/models"
)

type IReservationRepository interface {
	Create(item *models.Reservation) error
	GetUserReservationsCount(username string, status string) (int, error)
	GetUserReservations(username string) ([]*models.Reservation, error)
	GetReservation(ruid string) (*models.Reservation, error)
	UpdateReservation(uid string, status string) error
}
