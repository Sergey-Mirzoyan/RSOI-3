package repositories

import (
	"lab2/src/library/models"
)

type ILibraryRepository interface {
	GetByCity(city string) ([]*models.Library, error)
	GetByUid(luid string) (*models.Library, error)
}
