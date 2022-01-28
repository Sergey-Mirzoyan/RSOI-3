package usecase

import (
	"lab2/src/gateway/gateway_error"
	"lab2/src/gateway/models"
)

type IGatewayUsecase interface {
	GetCityLibraries(city string) (*[]models.Library, gateway_error.GatewayError)
	GetLibraryBooks(luid string, all bool) (*[]models.Book, gateway_error.GatewayError)
	GetUserReservations(username string) ([]*models.ReservationInfo, gateway_error.GatewayError)
	RentBook(username string, info *models.BookReservationInfo) (*models.ReservationResult, gateway_error.GatewayError)
	ReturnBook(reservationUuid string, info *models.BookReturningInfo) gateway_error.GatewayError
	GetUserRating(username string) (int, gateway_error.GatewayError)
}
