package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"lab2/src/rating/usecases"
	"net/http"
	"strconv"
)

type RatingHandlers struct {
	rc usecases.IRatingUsecase
}

func NewRatingHandlers(rating_case usecases.IRatingUsecase) *RatingHandlers {
	return &RatingHandlers{rc: rating_case}
}

func (rh *RatingHandlers) GetUserRating(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	username := r.Header.Get("X-User-Name")

	if username == "" {
		fmt.Println("Received an invalid header")
		http.Error(w, "Failed to get rating: incorrect username request header", http.StatusBadRequest)
		return
	}

	rating, err := rh.rc.GetByUser(username)
	if err != nil {
		fmt.Println("Failed to get user rating")
		http.Error(w, "Error while getting user rating", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rating)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (rh *RatingHandlers) PatchUserRating(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	username := r.Header.Get("X-User-Name")
	diff, err := strconv.Atoi(r.Header.Get("Rating-Difference"))

	if username == "" || err != nil {
		fmt.Println("Received an invalid header")
		http.Error(w, "Failed to get rating: incorrect request header", http.StatusBadRequest)
		return
	}

	err = rh.rc.AlterByUser(diff, username)
	if err != nil {
		fmt.Println("Failed to get user rating")
		http.Error(w, "Error while getting user rating", http.StatusInternalServerError)
		return
	}

	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}