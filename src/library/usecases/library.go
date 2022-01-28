package usecases

import "lab2/src/library/models"

type ILibraryUsecase interface {
	GetByCity(city string) ([]*models.Library, error)
	GetByUid(luid string) (*models.Library, error)
}
