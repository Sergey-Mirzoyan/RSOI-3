package usecases

import (
	"lab1/models"
	"lab1/repositories"
)

type IPersonsUsecase interface {
	GetAll() ([]*models.Person, error)
	GetById(id int) (*models.Person, error)
	Create(person *models.InputPerson) (int, error)
	Update(id int, newInfo *models.InputPerson) (models.Person, int, error)
	Delete(id int) (int, error)
}

type PersonsUsecase struct {
	pr repositories.IPersonsRepository
}

func NewPersonsUsecase(repo repositories.IPersonsRepository) *PersonsUsecase {
	return &PersonsUsecase{pr: repo}
}

func (pc *PersonsUsecase) Create(person *models.InputPerson) (int, error) {
	return pc.pr.Create(person)
}

func (pc *PersonsUsecase) GetById(id int) (*models.Person, error) {
	return pc.pr.GetById(id)
}

func (pc *PersonsUsecase) GetAll() ([]*models.Person, error) {
	return pc.pr.GetAll()
}

func (pc *PersonsUsecase) Update(id int, newInfo *models.InputPerson) (models.Person, int, error) {
	oldInfo, err := pc.pr.GetById(id)
	if err != nil {
		return models.Person{}, -1, err
	}
	if newInfo.Name == nil {
		newInfo.Name = &oldInfo.Name
	}
	if newInfo.Age == nil {
		newInfo.Age = &oldInfo.Age
	}
	if newInfo.Address == nil {
		newInfo.Address = &oldInfo.Address
	}
	if newInfo.Work == nil {
		newInfo.Work = &oldInfo.Work
	}
	affected, err := pc.pr.Update(id, newInfo)
	res := models.Person{
		Id:      id,
		Name:    *newInfo.Name,
		Age:     *newInfo.Age,
		Address: *newInfo.Address,
		Work:    *newInfo.Work,
	}
	return res, affected, err
}

func (pc *PersonsUsecase) Delete(id int) (int, error) {
	return pc.pr.Delete(id)
}
