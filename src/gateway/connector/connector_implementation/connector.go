package connector_implementation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lab2/src/gateway/connector"
	"lab2/src/gateway/models"
	"net/http"
	"strconv"
)

type GatewayConnector struct{
	config connector.Config
}

func NewGatewayConnector(config *connector.Config) *GatewayConnector {
	return &GatewayConnector{config: *config}
}

func (gc *GatewayConnector) GetCityLibraries(city string) (*[]models.Library, error) {
	url := fmt.Sprintf(gc.config.LibraryAddress + gc.config.ApiPath + "libraries?city=%s", city)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to get libraries from internal service")
		err = errors.New("Failed to get libraries from internal service")
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(response.Body)

	res := &[]models.Library{}
	err = json.NewDecoder(response.Body).Decode(res)
	if err != nil {
		fmt.Println("Failed to decode the json received from an internal service")
		err = errors.New("Internal decoding error")
		return nil, err
	}
	return res, nil
}

func (gc *GatewayConnector) GetLibraryBooks(luid string, all bool) (*[]models.Book, error) {
	url := fmt.Sprintf(gc.config.LibraryAddress + gc.config.ApiPath + "libraries/%s/books", luid)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return nil, err
	}

	if all {
		q := request.URL.Query()
		q.Add("shawAll", "true")
		request.URL.RawQuery = q.Encode()
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to get books from internal service")
		err = errors.New("Failed to get books from internal service")
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(response.Body)

	res := &[]models.Book{}
	err = json.NewDecoder(response.Body).Decode(res)
	if err != nil {
		fmt.Println("Failed to decode the json received from an internal service")
		err = errors.New("Internal decoding error")
		return nil, err
	}
	return res, nil
}

func (gc *GatewayConnector) GetUserReservationsCount(username string) (int, error) {
	url := fmt.Sprintf(gc.config.ReservationAddress + gc.config.ApiPath + "reservations/count")
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return -1, err
	}

	request.Header.Set("X-User-Name", username)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to get reservations from internal service")
		err = errors.New("Failed to get reservations from internal service")
		return -1, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(response.Body)

	var res int
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Println("Failed to decode the json received from an internal service")
		err = errors.New("Internal decoding error")
		return -1, err
	}
	return res, nil
}

func (gc *GatewayConnector) GetUserReservations(username string) ([]*models.Reservation, error) {
	url := fmt.Sprintf(gc.config.ReservationAddress + gc.config.ApiPath + "reservations")
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return nil, err
	}

	request.Header.Set("X-User-Name", username)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to get reservations from internal service")
		err = errors.New("Failed to get reservations from internal service")
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(response.Body)

	var res []*models.Reservation
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Println("Failed to decode the json received from an internal service")
		err = errors.New("Internal decoding error")
		return nil, err
	}
	return res, nil
}

func (gc *GatewayConnector) GetUserCurrentReservations(username string) (int, error) {
	url := fmt.Sprintf(gc.config.ReservationAddress + gc.config.ApiPath + "reservations/count")
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return -1, err
	}

	request.Header.Set("X-User-Name", username)
	request.Header.Set("Status", "RENTED")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to get reservations from internal service")
		err = errors.New("Failed to get reservations from internal service")
		return -1, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(response.Body)

	var res int
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Println("Failed to decode the json received from an internal service")
		err = errors.New("Internal decoding error")
		return -1, err
	}
	return res, nil
}

func (gc *GatewayConnector) GetUserRating(username string) (int, error) {
	url := fmt.Sprintf(gc.config.RatingAddress + gc.config.ApiPath + "rating")
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return -1, err
	}

	request.Header.Set("X-User-Name", username)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to get rating from internal service")
		err = errors.New("Failed to get rating from internal service")
		return -1, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(response.Body)

	var res int
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Println("Failed to decode the json received from an internal service")
		err = errors.New("Internal decoding error")
		return -1, err
	}
	return res, nil
}

func (gc *GatewayConnector) UpdateUserRating(username string, diff int) error {
	url := fmt.Sprintf(gc.config.RatingAddress + gc.config.ApiPath + "rating")
	request, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return err
	}

	request.Header.Set("X-User-Name", username)
	request.Header.Set("Rating-Difference", strconv.Itoa(diff))

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		fmt.Println("Failed to update rating in internal service")
		err = errors.New("Failed to update rating in internal service")
		return err
	}

	return nil
}

func (gc *GatewayConnector) UpdateBooksAmount(info *models.BookAmountInfo) error {
	url := fmt.Sprintf(gc.config.LibraryAddress + gc.config.ApiPath + "library_books")

	data, err := json.Marshal(info)
	if err != nil {
		fmt.Println("Failed to encode input data")
		err = errors.New("Encoding error")
		return err
	}
	request, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Failed to create an http request")
		return err
	}

	request.Header.Set("Library-Uid", info.LibraryUid)
	request.Header.Set("Book-Uid", info.BookUid)

	client := &http.Client{}
	_, err = client.Do(request)

	return err
}

func (gc *GatewayConnector) PostReservation(username string, info *models.BookReservationInfo) (string, error) {
	url := fmt.Sprintf(gc.config.ReservationAddress + gc.config.ApiPath + "reservations")
	data, err := json.Marshal(info)
	if err != nil {
		fmt.Println("Failed to encode input data")
		err = errors.New("Encoding error")
		return "", err
	}
	request, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Failed to create an http request")
		return "", err
	}

	request.Header.Set("X-User-Name", username)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to get reservationUid from internal service")
		err = errors.New("Failed to get reservationUid from internal service")
		return "", err
	}

	var uid string
	err = json.NewDecoder(response.Body).Decode(&uid)
	if err != nil {
		fmt.Println("Failed to decode data from internal service")
		err = errors.New("Decoding error")
		return "", err
	}

	return uid, err
}

func (gc *GatewayConnector) GetBooksCount(libraryUid string, bookUid string) (int, error) {
	url := fmt.Sprintf(gc.config.LibraryAddress + gc.config.ApiPath + "library_books")
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return -1, err
	}

	request.Header.Set("Library-Uid", libraryUid)
	request.Header.Set("Book-Uid", bookUid)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to get books count from internal service")
		err = errors.New("Failed to get books count from internal service")
		return -1, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(response.Body)

	var res int
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		fmt.Println("Failed to decode the json received from an internal service")
		err = errors.New("Internal decoding error")
		return -1, err
	}
	return res, nil
}

func (gc *GatewayConnector) GetReservation(reservationUid string) (*models.Reservation, error) {
	url := fmt.Sprintf(gc.config.ReservationAddress + gc.config.ApiPath + "reservations/%s", reservationUid)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to get reservation from internal service")
		err = errors.New("Failed to get reservation from internal service")
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(response.Body)

	res := &models.Reservation{}
	err = json.NewDecoder(response.Body).Decode(res)
	if err != nil {
		fmt.Println("Failed to decode the json received from an internal service")
		err = errors.New("Internal decoding error")
		return nil, err
	}
	return res, nil
}

func (gc *GatewayConnector) GetBook(bookUid string) (*models.Book, error) {
	url := fmt.Sprintf(gc.config.LibraryAddress + gc.config.ApiPath + "books/%s", bookUid)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to get a book from internal service")
		err = errors.New("Failed to get the book from internal service")
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(response.Body)

	res := &models.Book{}
	err = json.NewDecoder(response.Body).Decode(res)
	if err != nil {
		fmt.Println("Failed to decode the json received from an internal service")
		err = errors.New("Internal decoding error")
		return nil, err
	}
	return res, nil
}

func (gc *GatewayConnector) UpdateBookCondition(bookUid string, condition string) error {
	url := fmt.Sprintf(gc.config.LibraryAddress + gc.config.ApiPath + "books/%s", bookUid)
	info := models.BookPatch{ Condition: &condition }
	data, err := json.Marshal(info)
	if err != nil {
		fmt.Println("Failed to encode input data")
		err = errors.New("Encoding error")
		return err
	}
	request, err := http.NewRequest("PUT", url, bytes.NewReader(data))

	client := &http.Client{}
	_, err = client.Do(request)

	return err
}

func (gc *GatewayConnector) UpdateReservationStatus(reservationUid string, status string) error {
	url := fmt.Sprintf(gc.config.ReservationAddress + gc.config.ApiPath + "reservations/%s", reservationUid)
	info := models.ReservationPatch{Status: &status}
	data, err := json.Marshal(info)
	if err != nil {
		fmt.Println("Failed to encode input data")
		err = errors.New("Encoding error")
		return err
	}
	request, err := http.NewRequest("PUT", url, bytes.NewReader(data))

	client := &http.Client{}
	_, err = client.Do(request)

	return err
}

func (gc *GatewayConnector) GetLibrary(libraryUid string) (*models.Library, error) {
	url := fmt.Sprintf(gc.config.LibraryAddress + gc.config.ApiPath + "libraries/%s", libraryUid)
	response, err := http.Get(url)

	res := &models.Library{}

	err = json.NewDecoder(response.Body).Decode(res)
	if err != nil {
		fmt.Println("Failed to decode the json received from an internal service")
		err = errors.New("Internal decoding error")
		return nil, err
	}

	return res, err
}