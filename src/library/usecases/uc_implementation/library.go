package uc_implementation

import (
	"lab2/src/library/models"
	"lab2/src/library/repositories"
)

type LibraryUsecase struct {
	lr repositories.ILibraryRepository
}

func NewLibraryUsecase(repo repositories.ILibraryRepository) *LibraryUsecase {
	return &LibraryUsecase{lr: repo}
}

func (lc *LibraryUsecase) GetByUid(luid string) (*models.Library, error) {
	return lc.lr.GetByUid(luid)
}

func (lc *LibraryUsecase) GetByCity(city string) ([]*models.Library, error) {
	return lc.lr.GetByCity(city)
}
