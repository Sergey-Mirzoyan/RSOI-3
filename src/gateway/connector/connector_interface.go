package connector

import "lab2/src/gateway/models"

type IGatewayConnector interface {
	GetCityLibraries(city string) (*[]models.Library, error)
	GetLibraryBooks(luid string, all bool) (*[]models.Book, error)
	GetUserReservations(username string) ([]*models.Reservation, error)
	GetUserReservationsCount(username string) (int, error)
	GetUserCurrentReservations(username string) (int, error)
	GetUserRating(username string) (int, error)
	PostReservation(username string, info *models.BookReservationInfo) (string, error)
	GetBooksCount(libraryUid string, bookUid string) (int, error)
	GetReservation(reservationUid string) (*models.Reservation, error)
	GetBook(bookUid string) (*models.Book, error)
	UpdateBookCondition(bookUid string, condition string) error
	UpdateBooksAmount(info *models.BookAmountInfo) error
	GetLibrary(libraryUid string) (*models.Library, error)
	UpdateReservationStatus(reservationUid string, status string) error
	UpdateUserRating(username string, diff int) error
}
