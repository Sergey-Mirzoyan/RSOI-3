package uc_implementation

import (
	"errors"
	"fmt"
	"lab2/src/gateway/connector"
	"lab2/src/gateway/gateway_error"
	"lab2/src/gateway/models"
	"lab2/src/gateway/request_queue"
	"lab2/src/gateway/usecase"
	"time"
)

type GatewayUsecase struct {
	connector     connector.IGatewayConnector
	libraryCB     usecase.CircuitBreaker
	reservationCB usecase.CircuitBreaker
	ratingCB      usecase.CircuitBreaker
	libraryQueue  request_queue.QueueRepeater
	ratingQueue   request_queue.QueueRepeater
}

func NewGatewayUsecase(connector connector.IGatewayConnector) *GatewayUsecase {
	uc := &GatewayUsecase{
		connector:     connector,
		libraryCB:     *usecase.NewCircuitBreaker(50),
		reservationCB: *usecase.NewCircuitBreaker(50),
		ratingCB:      *usecase.NewCircuitBreaker(50),
		libraryQueue:  *request_queue.NewQueueRepeater(),
		ratingQueue:   *request_queue.NewQueueRepeater(),
	}
	uc.libraryQueue.Start()
	uc.ratingQueue.Start()
	return uc
}

func (gc *GatewayUsecase) Close() {
	gc.libraryQueue.Stop()
	gc.ratingQueue.Stop()
}

func (gc *GatewayUsecase) GetCityLibraries(city string) (*[]models.Library, gateway_error.GatewayError) {
	code := gateway_error.Ok
	res, err := gc.libraryCB.Call(func() (interface{}, error) { return gc.connector.GetCityLibraries(city) })
	if err != nil {
		code = gateway_error.Internal
	}
	return res.(*[]models.Library), gateway_error.GatewayError{Err: err, Code: code}
}

func (gc *GatewayUsecase) GetLibraryBooks(luid string, all bool) (*[]models.Book, gateway_error.GatewayError) {
	code := gateway_error.Ok
	res, err := gc.libraryCB.Call(func() (interface{}, error) { return gc.connector.GetLibraryBooks(luid, all) })
	if err != nil {
		code = gateway_error.Internal
	}
	return res.(*[]models.Book), gateway_error.GatewayError{Err: err, Code: code}
}

func (gc *GatewayUsecase) GetUserReservations(username string) ([]*models.ReservationInfo, gateway_error.GatewayError) {
	code := gateway_error.Ok
	res, err := gc.reservationCB.Call(func() (interface{}, error) { return gc.connector.GetUserReservations(username) })
	if err != nil {
		code = gateway_error.Internal
	}

	info := make([]*models.ReservationInfo, 0)
	for _, reservation := range res.([]*models.Reservation) {
		book, err := gc.connector.GetBook(reservation.BookUid)
		var shortBook *models.ShortBook
		if err != nil {
			shortBook = &models.ShortBook{
				BookUid: reservation.BookUid,
			}
		} else {
			shortBook = &models.ShortBook{
				BookUid:   book.BookUid,
				Name:      book.Name,
				Author:    book.Author,
				Genre:     book.Genre,
				Condition: book.Condition,
			}
		}
		library, err := gc.connector.GetLibrary(reservation.LibraryUid)
		var shortLibrary *models.ShortLibrary
		if err != nil {
			shortLibrary = &models.ShortLibrary{
				LibraryUid: reservation.LibraryUid,
			}
		} else {
			shortLibrary = &models.ShortLibrary{
				LibraryUid: library.LibraryUid,
				Name:       library.Name,
				City:       library.City,
				Address:    library.Address,
			}
		}
		info = append(info, &models.ReservationInfo{
			ReservationUid: reservation.ReservationUid,
			Status:         reservation.Status,
			StartDate:      reservation.StartDate,
			TillDate:       reservation.TillDate,
			Book:           shortBook,
			Library:        shortLibrary,
		})
	}

	return info, gateway_error.GatewayError{Err: err, Code: code}
}

func (gc *GatewayUsecase) RentBook(username string, info *models.BookReservationInfo) (*models.ReservationResult, gateway_error.GatewayError) {
	reserved, err := gc.connector.GetUserCurrentReservations(username)
	if err != nil {
		fmt.Printf("Failed to get user reservations\n")
		return nil, gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
	}

	allowed, err := gc.connector.GetUserRating(username)
	if err != nil {
		fmt.Printf("Failed to get user rating\n")
		return nil, gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
	}

	if reserved < allowed {
		count, err := gc.connector.GetBooksCount(info.LibraryUid, info.BookUid)
		if err != nil {
			fmt.Printf("Failed to get books amount\n")
			return nil, gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
		}

		if count == 0 {
			fmt.Printf("No books left for rent\n")
			err = errors.New("No books left for rent")
			return nil, gateway_error.GatewayError{Err: err, Code: gateway_error.User}
		}

		amountInfo := &models.BookAmountInfo{
			BookUid:        info.BookUid,
			LibraryUid:     info.LibraryUid,
			AvailableCount: count - 1,
		}
		err = gc.connector.UpdateBooksAmount(amountInfo)
		if err != nil {
			fmt.Printf("Failed to occupy a book\n")
			return nil, gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
		}

		uid, err := gc.connector.PostReservation(username, info)
		if err != nil {
			fmt.Printf("Failed to rent a book\n")
			return nil, gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
		}

		book, err := gc.connector.GetBook(info.BookUid)
		if err != nil {
			fmt.Printf("Failed to get a book\n")
			return nil, gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
		}

		library, err := gc.connector.GetLibrary(info.LibraryUid)
		if err != nil {
			fmt.Printf("Failed to get a library\n")
			return nil, gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
		}

		rating, err := gc.connector.GetUserRating(username)
		if err != nil {
			fmt.Printf("Failed to get user rating\n")
			return nil, gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
		}

		reservationInfo := &models.ReservationResult{
			ReservationUid: uid,
			Status:         "RENTED",
			StartDate:      time.Now().Format("2006-01-02"),
			TillDate:       info.TillDate,
			Book: &models.ShortBook{
				BookUid:   book.BookUid,
				Name:      book.Name,
				Author:    book.Author,
				Genre:     book.Genre,
				Condition: book.Condition,
			},
			Library: &models.ShortLibrary{
				LibraryUid: library.LibraryUid,
				Name:       library.Name,
				City:       library.City,
				Address:    library.Address,
			},
			Rating: &models.Rating{Stars: rating},
		}

		return reservationInfo, gateway_error.GatewayError{Err: nil, Code: gateway_error.Ok}
	} else {
		fmt.Printf("Rent failed: user has already reached their reservation limit\n")
		err = errors.New("The user has already reached their reservation limit")
		return nil, gateway_error.GatewayError{Err: err, Code: gateway_error.User}
	}
}

func (gc *GatewayUsecase) ReturnBook(reservationUuid string, info *models.BookReturningInfo) gateway_error.GatewayError {
	diff := 0
	reservation, err := gc.connector.GetReservation(reservationUuid)
	if err != nil {
		fmt.Printf("Failed to get reservation")
		err = errors.New("Invalid reservation ID")
		return gateway_error.GatewayError{Err: err, Code: gateway_error.User}
	}
	book, err := gc.connector.GetBook(reservation.BookUid)
	if err != nil {
		fmt.Printf("Failed to get book")
		err = errors.New("Failed to get reserved book data")
		return gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
	}
	if info.Condition != book.Condition {
		diff -= 10
	}
	date, err := time.Parse("2006-01-02", info.Date)
	if err != nil {
		fmt.Printf("Failed to decode reservation end date")
		err = errors.New("Failed to decode reservation end date")
		return gateway_error.GatewayError{Err: err, Code: gateway_error.User}
	}
	till, err := time.Parse("2006-01-02", reservation.TillDate)
	if err != nil {
		fmt.Printf("Failed to decode reservation limit date")
		err = errors.New("Failed to decode reservation limit date")
		return gateway_error.GatewayError{Err: err, Code: gateway_error.User}
	}
	status := "RETURNED"
	if date.After(till) {
		status = "EXPIRED"
		diff -= 10
	}
	if diff == 0 {
		diff = 1
	}

	err = gc.connector.UpdateReservationStatus(reservationUuid, status)
	if err != nil {
		fmt.Printf("Failed to update reservation status")
		err = errors.New("Failed to update reservation status")
		return gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
	}

	err = gc.connector.UpdateUserRating(reservation.Username, diff)
	if err != nil {
		gc.ratingQueue.AddRequest(func() error { return gc.connector.UpdateUserRating(reservation.Username, diff) }, 10*time.Second)
	}

	if info.Condition != book.Condition {
		err = gc.connector.UpdateBookCondition(reservation.BookUid, info.Condition)
		if err != nil {
			fmt.Printf("Failed to update book condition")
			err = errors.New("Failed to update book condition")
			return gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
		}
	}

	count, err := gc.connector.GetBooksCount(reservation.LibraryUid, reservation.BookUid)
	if err != nil {
		fmt.Printf("Failed to get books amount\n")
		return gateway_error.GatewayError{Err: err, Code: gateway_error.Internal}
	}
	amountInfo := &models.BookAmountInfo{
		BookUid:        reservation.BookUid,
		LibraryUid:     reservation.LibraryUid,
		AvailableCount: count + 1,
	}
	err = gc.connector.UpdateBooksAmount(amountInfo)
	if err != nil {
		gc.ratingQueue.AddRequest(func() error { return gc.connector.UpdateBooksAmount(amountInfo) }, 10*time.Second)
	}
	return gateway_error.GatewayError{Err: nil, Code: gateway_error.Ok}
}

func (gc *GatewayUsecase) GetUserRating(username string) (int, gateway_error.GatewayError) {
	code := gateway_error.Ok
	res, err := gc.ratingCB.Call(func() (interface{}, error) { return gc.connector.GetUserRating(username) })
	if err != nil {
		code = gateway_error.Internal
	}
	return res.(int), gateway_error.GatewayError{Err: err, Code: code}
}
