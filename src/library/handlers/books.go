package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"lab2/src/library/models"
	"lab2/src/library/usecases"
	"net/http"
)

type BooksHandlers struct {
	bc usecases.IBooksUsecase
}

func NewBooksHandlers(booksCase usecases.IBooksUsecase) *BooksHandlers {
	return &BooksHandlers{bc: booksCase}
}

func (bh *BooksHandlers) GetBook(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	vars := mux.Vars(r)
	luid := vars["bookUid"]
	if luid == "" {
		fmt.Println("Received an invalid path parameter")
		http.Error(w, "Failed to find libraries: wrong path parameter", http.StatusBadRequest)
		return
	}

	foundBook, err := bh.bc.GetByUid(luid)
	if err != nil {
		fmt.Println("Failed to get books")
		http.Error(w, "Error while getting books list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(foundBook)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (bh *BooksHandlers) GetByLibrary(w http.ResponseWriter, r *http.Request) {
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

	all := false
	query := r.URL.Query()
	showAll, present := query["showAll"]
	if present {
		if len(showAll) != 1 {
			fmt.Println("Received an invalid query parameter")
			http.Error(w, "Failed to find libraries: wrong query parameter", http.StatusBadRequest)
			return
		}
		if showAll[0] == "true" {
			all = true
		}
	}

	foundBooks, err := bh.bc.GetByLibraryUid(luid, all)
	if err != nil {
		fmt.Println("Failed to get books")
		http.Error(w, "Error while getting books list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(foundBooks)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (bh *BooksHandlers) UpdateBook(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	vars := mux.Vars(r)
	buid := vars["bookUid"]
	if buid == "" {
		fmt.Println("Received an invalid path parameter")
		http.Error(w, "Failed to find libraries: wrong path parameter", http.StatusBadRequest)
		return
	}

	info := &models.BookPatch{}
	err := json.NewDecoder(r.Body).Decode(info)
	if err != nil {
		fmt.Println("Decoding json error: ", err)
		http.Error(w, "Failed to decode json data", http.StatusBadRequest)
		return
	}

	err = bh.bc.UpdateBook(buid, info)
	if err != nil {
		fmt.Println("Failed to get books")
		http.Error(w, "Error while getting books list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

