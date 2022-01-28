package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"lab2/src/library/models"
	"lab2/src/library/usecases"
	"net/http"
)

type LibraryBooksHandlers struct {
	lbc usecases.ILibraryBooksUsecase
}

func NewLibraryBooksHandlers(booksCase usecases.ILibraryBooksUsecase) *LibraryBooksHandlers {
	return &LibraryBooksHandlers{lbc: booksCase}
}

func (lbh *LibraryBooksHandlers) UpdateBooksAmount(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	info := &models.LibraryBooks{}
	err := json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		fmt.Println("Failed to decode the received json")
		http.Error(w, "Bad json given as input", http.StatusBadRequest)
		return
	}

	luid := r.Header.Get("Library-Uid")
	buid := r.Header.Get("Book-Uid")
	if luid == "" || buid == "" {
		fmt.Println("Received an invalid header")
		http.Error(w, "Failed to get reservations: incorrect library/book uid combination in headers", http.StatusBadRequest)
		return
	}

	err = lbh.lbc.UpdateBooksAmount(luid, buid, info.AvailableCount)
	if err != nil {
		fmt.Println("Failed to decrease books amount")
		http.Error(w, "Error while decreasing books amount list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (lbh *LibraryBooksHandlers) GetBooksAmount(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	luid := r.Header.Get("Library-Uid")
	buid := r.Header.Get("Book-Uid")
	if luid == "" || buid == "" {
		fmt.Println("Received an invalid header")
		http.Error(w, "Failed to get reservations: incorrect library/book uid combination in headers", http.StatusBadRequest)
		return
	}

	res, err := lbh.lbc.GetBooksAmount(luid, buid)
	if err != nil {
		fmt.Println("Failed to decrease books amount")
		http.Error(w, "Error while decreasing books amount list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}