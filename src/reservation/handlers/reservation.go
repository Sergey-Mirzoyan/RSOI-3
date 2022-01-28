package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"lab2/src/reservation/models"
	"lab2/src/reservation/usecases"
	"net/http"
)

type ReservationHandlers struct {
	rc usecases.IReservationUsecase
}

func NewReservationHandlers(reservation_case usecases.IReservationUsecase) *ReservationHandlers {
	return &ReservationHandlers{rc: reservation_case}
}

func (rh *ReservationHandlers) GetReservation(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)


	vars := mux.Vars(r)
	ruid := vars["reservationUid"]
	if ruid == "" {
		fmt.Println("Received an invalid path parameter")
		http.Error(w, "Failed to find libraries: wrong path parameter", http.StatusBadRequest)
		return
	}

	reservations, err := rh.rc.GetReservation(ruid)
	if err != nil {
		fmt.Println("Failed to get user reservations")
		http.Error(w, "Error while getting user reservations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reservations)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (rh *ReservationHandlers) GetUserReservationsCount(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	username := r.Header.Get("X-User-Name")
	status := r.Header.Get("Status")

	if username == "" {
		fmt.Println("Received an invalid header")
		http.Error(w, "Failed to get reservations: incorrect username request header", http.StatusBadRequest)
		return
	}

	reservations, err := rh.rc.GetUserReservationsCount(username, status)
	if err != nil {
		fmt.Println("Failed to get user reservations")
		http.Error(w, "Error while getting user reservations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reservations)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (rh *ReservationHandlers) GetUserReservations(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	username := r.Header.Get("X-User-Name")

	if username == "" {
		fmt.Println("Received an invalid header")
		http.Error(w, "Failed to get reservations: incorrect username request header", http.StatusBadRequest)
		return
	}

	reservations, err := rh.rc.GetUserReservations(username)
	if err != nil {
		fmt.Println("Failed to get user reservations")
		http.Error(w, "Error while getting user reservations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reservations)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (rh *ReservationHandlers) PostReservation(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	newReservation := &models.InputReservation{}
	err := json.NewDecoder(r.Body).Decode(newReservation)
	if err != nil {
		fmt.Println("Failed to decode the received json")
		http.Error(w, "Bad json given as input", http.StatusBadRequest)
		return
	}

	username := r.Header.Get("X-User-Name")
	if username == "" {
		fmt.Println("Received an invalid header")
		http.Error(w, "Failed to get reservations: incorrect username request header", http.StatusBadRequest)
		return
	}

	newReservation.Username = username
	uid, err := rh.rc.Create(newReservation)
	if err != nil {
		fmt.Println("Failed to create a new person record")
		http.Error(w, "Failed to create a person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(uid)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (rh *ReservationHandlers) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	patch := &models.ReservationPatch{}
	err := json.NewDecoder(r.Body).Decode(patch)
	if err != nil {
		fmt.Println("Failed to decode the received json")
		http.Error(w, "Bad json given as input", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	uid := vars["reservationUid"]
	if uid == "" {
		fmt.Println("Received an invalid path parameter")
		http.Error(w, "Failed to find libraries: wrong path parameter", http.StatusBadRequest)
		return
	}

	err = rh.rc.UpdateReservation(uid, *patch.Status)
	if err != nil {
		fmt.Println("Failed to create a new person record")
		http.Error(w, "Failed to create a person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}