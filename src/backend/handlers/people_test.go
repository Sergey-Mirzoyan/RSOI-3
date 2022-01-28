package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"lab1/mock"
	"lab1/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

type TestCaseCreate struct {
	Person    *models.InputPerson
	PersonID  int
	Response  string
	InputErr  bool
	CreateErr bool
	Status    int
}

func TestPersonsHandlers_Create(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	tuc := mock.NewMockIPersonsUsecase(ctl)
	handler := PersonsHandlers{
		pc: tuc,
	}

	names := []string{"name", "", "name"}
	ages := []int{12, 0, 35}
	addresses := []string{"address", "", "addr"}
	works := []string{"he is only twelve", "", "example"}

	cases := []TestCaseCreate{
		{
			Person:    &models.InputPerson{},
			PersonID:  1,
			Response:  "",
			InputErr:  false,
			CreateErr: false,
			Status:    201,
		},
		{
			Person:    &models.InputPerson{},
			PersonID:  1,
			Response:  "Bad json given as input\n",
			InputErr:  true,
			CreateErr: false,
			Status:    400,
		},
		{
			Person:    &models.InputPerson{},
			PersonID:  1,
			Response:  "Failed to create a person\n",
			InputErr:  false,
			CreateErr: true,
			Status:    500,
		},
	}

	for i, _ := range cases {
		if !cases[i].InputErr {
			cases[i].Person.Name = &names[i]
			cases[i].Person.Age = &ages[i]
			cases[i].Person.Address = &addresses[i]
			cases[i].Person.Work = &works[i]
		}
	}

	for caseNum, item := range cases {
		var r *http.Request

		if item.InputErr == false {
			jsonPerson, _ := json.Marshal(item.Person)
			r = httptest.NewRequest("POST", "/api/pins", strings.NewReader(string(jsonPerson)))
		} else {
			r = httptest.NewRequest("POST", "/api/pins", strings.NewReader("blabla"))
		}
		r.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		if !item.InputErr {
			err := error(nil)
			id := item.PersonID
			if item.CreateErr {
				err = errors.New("")
				id = 0
			}
			tuc.EXPECT().Create(item.Person).Return(id, err)
		}

		handler.PostPerson(w, r)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)

		if item.CreateErr {
			if w.Code == 200 || w.Code == 201 {
				t.Errorf("[%d] wrong status Response: got %+v, expected not success status",
					caseNum, w.Code)
			}
		} else {
			if bodyStr != item.Response || w.Code != item.Status {
				t.Errorf("[%d] wrong Response: got %+v, code: %d;\nexpected %+v, code: %d",
					caseNum, bodyStr, w.Code, item.Response, item.Status)
			}
		}
	}
}

type TestCaseGet struct {
	Person   *models.Person
	PersonID int
	Response string
	IdErr    bool
	InputErr bool
	Status   int
}

func TestPersonsHandlers_GetPersonById(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	tuc := mock.NewMockIPersonsUsecase(ctl)
	handler := PersonsHandlers{
		pc: tuc,
	}

	validPerson := models.Person{
		Id:      1,
		Name:    "name",
		Age:     36,
		Address: "address",
		Work:    "example",
	}
	validJSON, _ := json.Marshal(validPerson)

	cases := []TestCaseGet{
		{
			Person:   &validPerson,
			PersonID: 1,
			Response: string(validJSON) + "\n",
			InputErr: false,
			IdErr:    false,
			Status:   200,
		},
		{
			Person:   &validPerson,
			PersonID: -1,
			Response: "No person matches the provided id\n",
			InputErr: false,
			IdErr:    true,
			Status:   404,
		},
		{
			Person:   &validPerson,
			PersonID: 1,
			Response: "Failed to find a person: wrong path parameter\n",
			IdErr:    false,
			InputErr: true,
			Status:   400,
		},
	}

	for caseNum, item := range cases {
		var r *http.Request

		if item.IdErr == false {
			jsonPerson, _ := json.Marshal(item.Person)
			r = httptest.NewRequest("POST", "/api/pins", strings.NewReader(string(jsonPerson)))
		} else {
			r = httptest.NewRequest("POST", "/api/pins", strings.NewReader("blabla"))
		}
		r.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		if !item.InputErr {
			if !item.IdErr {
				tuc.EXPECT().GetById(item.PersonID).Return(item.Person, error(nil))
			} else {
				tuc.EXPECT().GetById(item.PersonID).Return(nil, errors.New(""))
			}
		}

		if !item.InputErr {
			r = mux.SetURLVars(r, map[string]string{"id": fmt.Sprintf("%d", item.PersonID)})
		}

		handler.GetPersonById(w, r)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)

		if bodyStr != item.Response || w.Code != item.Status {
			t.Errorf("[%d] wrong Response: got %+v, code: %d;\nexpected %+v, code: %d",
				caseNum, bodyStr, w.Code, item.Response, item.Status)
		}
	}
}
