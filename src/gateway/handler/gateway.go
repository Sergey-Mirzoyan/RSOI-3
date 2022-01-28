package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"lab2/src/gateway/models"
	"lab2/src/gateway/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type GatewayHandlers struct {
	gc usecase.IGatewayUsecase
}

func NewGatewayHandlers(gateway_case usecase.IGatewayUsecase) *GatewayHandlers {
	return &GatewayHandlers{gc: gateway_case}
}

func writeError(w http.ResponseWriter, err string, status int) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err)))
}

func (gh *GatewayHandlers) GetCityLibraries(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	city, present := query["city"]
	if !present || len(city) != 1 {
		fmt.Println("Received an invalid query parameter")
		writeError(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
		return
	}

	foundLibraries, ret := gh.gc.GetCityLibraries(city[0])
	if ret.Err != nil {
		fmt.Println("Failed to get libraries")
		writeError(w, "Error while getting libraries list", http.StatusServiceUnavailable)
		return
	}

	param, present := query["page"]
	if !present || len(param) != 1 {
		fmt.Println("Received an invalid query parameter")
		writeError(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(param[0])
	if err != nil {
		fmt.Println("Received an invalid query parameter")
		writeError(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
		return
	}

	param, present = query["size"]
	if !present || len(param) != 1 {
		fmt.Println("Received an invalid query parameter")
		writeError(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
		return
	}
	pageSize, err := strconv.Atoi(param[0])
	if err != nil {
		fmt.Println("Received an invalid query parameter")
		writeError(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
		return
	}

	fin := page * pageSize
	if len(*foundLibraries) < fin {
		fin = len(*foundLibraries)
	}
	pageLibraries := (*foundLibraries)[(page-1)*pageSize : fin]
	res := models.LibrariesPage{
		Page:          page,
		PageSize:      pageSize,
		TotalElements: len(pageLibraries),
		Items:         &pageLibraries,
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		writeError(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (gh *GatewayHandlers) GetLibraryBooks(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	luid := vars["libraryUid"]
	if luid == "" {
		fmt.Println("Received an invalid path parameter")
		writeError(w, "Failed to find libraries: wrong path parameter", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	all := false
	showAll, present := query["showAll"]
	if !present || len(showAll) != 1 {
		fmt.Println("Received an invalid query parameter")
		writeError(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
		return
	}
	if showAll[0] == "true" {
		all = true
	}

	foundBooks, ret := gh.gc.GetLibraryBooks(luid, all)
	if ret.Err != nil {
		fmt.Println("Failed to get books")
		writeError(w, "Error while getting books list", http.StatusServiceUnavailable)
		return
	}

	param, present := query["page"]
	if !present || len(param) != 1 {
		fmt.Println("Received an invalid query parameter")
		writeError(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(param[0])
	if err != nil {
		fmt.Println("Received an invalid query parameter")
		writeError(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
		return
	}

	param, present = query["size"]
	if !present || len(param) != 1 {
		fmt.Println("Received an invalid query parameter")
		writeError(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
		return
	}
	pageSize, err := strconv.Atoi(param[0])
	if err != nil {
		fmt.Println("Received an invalid query parameter")
		writeError(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
		return
	}

	fin := page * pageSize
	if len(*foundBooks) < fin {
		fin = len(*foundBooks)
	}
	pageBooks := (*foundBooks)[(page-1)*pageSize : fin]
	res := models.BooksPage{
		Page:          page,
		PageSize:      pageSize,
		TotalElements: len(pageBooks),
		Items:         &pageBooks,
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		writeError(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (gh *GatewayHandlers) GetUserReservations(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	w.Header().Set("Content-Type", "application/json")

	username := r.Header.Get("X-User-Name")

	if username == "" {
		fmt.Println("Received an invalid header")
		writeError(w, "Failed to get reservations: incorrect username request header", http.StatusBadRequest)
		return
	}

	reservations, res := gh.gc.GetUserReservations(username)
	if res.Err != nil {
		fmt.Println("Failed to get user reservations")
		writeError(w, "Error while getting user reservations", http.StatusServiceUnavailable)
		return
	}

	err := json.NewEncoder(w).Encode(reservations)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		writeError(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (gh *GatewayHandlers) RentBook(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	w.Header().Set("Content-Type", "application/json")

	username := r.Header.Get("X-User-Name")

	if username == "" {
		fmt.Println("Received an invalid header")
		writeError(w, "Failed to get reservations: incorrect username request header", http.StatusBadRequest)
		return
	}

	info := &models.BookReservationInfo{}
	err := json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		fmt.Println("Failed to decode the received json")
		writeError(w, "Bad json given as input", http.StatusBadRequest)
		return
	}

	reservationInfo, res := gh.gc.RentBook(username, info)
	if res.Err != nil {
		fmt.Println("Failed to rent a book")
		writeError(w, "Bonus Service unavailable", http.StatusServiceUnavailable)
		return
	}

	err = json.NewEncoder(w).Encode(&reservationInfo)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		writeError(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (gh *GatewayHandlers) ReturnBook(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	time.Sleep(2 * time.Second)

	vars := mux.Vars(r)
	ruid := vars["reservationUid"]
	if ruid == "" {
		fmt.Println("Received an invalid path parameter")
		writeError(w, "Failed to return the book: wrong path parameter", http.StatusBadRequest)
		return
	}

	info := &models.BookReturningInfo{}
	err := json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		fmt.Println("Failed to decode the received json")
		writeError(w, "Bad json given as input", http.StatusBadRequest)
		return
	}

	res := gh.gc.ReturnBook(ruid, info)
	if res.Err != nil {
		fmt.Println("Failed to return a book")
		writeError(w, "Failed to return the book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (gh *GatewayHandlers) GetUserRating(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	time.Sleep(2 * time.Second)

	w.Header().Set("Content-Type", "application/json")

	username := r.Header.Get("X-User-Name")

	if username == "" {
		fmt.Println("Received an invalid header")
		writeError(w, "Failed to get rating: incorrect username request header", http.StatusBadRequest)
		return
	}

	rating, res := gh.gc.GetUserRating(username)
	if res.Err != nil {
		fmt.Println("Failed to get user rating")
		writeError(w, "Bonus Service unavailable", http.StatusServiceUnavailable)
		return
	}

	err := json.NewEncoder(w).Encode(&models.Rating{Stars: rating})
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		writeError(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}
