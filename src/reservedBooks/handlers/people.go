package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"lab1/models"
	"lab1/usecases"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PersonsHandlers struct {
	pc usecases.IPersonsUsecase
}

func NewPerosnsHandlers(person_case usecases.IPersonsUsecase) *PersonsHandlers {
	return &PersonsHandlers{pc: person_case}
}

func (ph *PersonsHandlers) PostPerson(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	newPerson := &models.InputPerson{}
	err := json.NewDecoder(r.Body).Decode(newPerson)
	if err != nil {
		fmt.Println("Failed to decode the received json")
		http.Error(w, "Bad json given as input", http.StatusBadRequest)
		return
	}

	id, err := ph.pc.Create(newPerson)
	if err != nil {
		fmt.Println("Failed to create a new person record")
		http.Error(w, "Failed to create a person", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/api/v1/persons/%d", id))
	w.WriteHeader(http.StatusCreated)
}

func (ph *PersonsHandlers) GetPersonById(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("Received an invalid path parameter")
		http.Error(w, "Failed to find a person: wrong path parameter", http.StatusBadRequest)
		return
	}

	foundPerson, err := ph.pc.GetById(id)
	if err != nil {
		fmt.Println("Person not found")
		http.Error(w, "No person matches the provided id", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(foundPerson)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (ph *PersonsHandlers) GetPersons(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	foundPeople, err := ph.pc.GetAll()
	if err != nil {
		fmt.Println("Failed to get people")
		http.Error(w, "Error while getting people list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(foundPeople)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (ph *PersonsHandlers) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("Received an invalid path parameter")
		http.Error(w, "Failed to find a person: wrong path parameter", http.StatusBadRequest)
		return
	}

	person := &models.InputPerson{}
	err = json.NewDecoder(r.Body).Decode(person)
	if err != nil {
		fmt.Println("Failed to decode the received json")
		http.Error(w, "Bad json given as input", http.StatusBadRequest)
		return
	}

	updatedPerson, affectedRows, err := ph.pc.Update(id, person)
	if err != nil {
		if affectedRows == 0 {
			fmt.Println("Bad id given to update a person")
			http.Error(w, "No person that matches the given id was found in database", http.StatusNotFound)
		} else {
			fmt.Println("Failed to update a person")
			http.Error(w, "Error while updating the person", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedPerson)
	if err != nil {
		fmt.Println("Encoding json error: ", err)
		http.Error(w, "Failed to encode data to json", http.StatusInternalServerError)
		return
	}
}

func (ph *PersonsHandlers) DeletePerson(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(r.Body)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("Received an invalid path parameter")
		http.Error(w, "Failed to find a person: wrong path parameter", http.StatusBadRequest)
		return
	}

	affectedRows, err := ph.pc.Delete(id)
	if err != nil {
		if affectedRows == 0 {
			fmt.Println("Bad id given to delete a person")
			http.Error(w, "No person that matches the given id was found in database", http.StatusNotFound)
		} else {
			fmt.Println("Failed to delete a person")
			http.Error(w, "Error while updating the person", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
