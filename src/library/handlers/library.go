package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"lab2/src/library/usecases"
	"net/http"
)

type LibraryHandlers struct {
	lc usecases.ILibraryUsecase
}

func NewLibraryHandlers(library_case usecases.ILibraryUsecase) *LibraryHandlers {
	return &LibraryHandlers{lc: library_case}
}

func (lh *LibraryHandlers) GetByUid(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	vars := mux.Vars(r)
	luid := vars["libraryUid"]
	if luid == "" {
		fmt.Println("Received an invalid path parameter")
		http.Error(w, "Failed to find libraries: wrong path parameter", http.StatusBadRequest)
		return
	}

	foundLibrary, err := lh.lc.GetByUid(luid)
	if err != nil {
		fmt.Println("Failed to get libraries")
		http.Error(w, "Error while getting libraries list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(foundLibrary)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (lh *LibraryHandlers) GetByCity(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	query := r.URL.Query()
	city, present := query["city"]
	if !present || len(city) != 1 {
		fmt.Println("Received an invalid query parameter")
		http.Error(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
		return
	}

	foundLibraries, err := lh.lc.GetByCity(city[0])
	if err != nil {
		fmt.Println("Failed to get libraries")
		http.Error(w, "Error while getting libraries list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(foundLibraries)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}
