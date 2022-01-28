package usecases

import "lab2/src/library/models"

type IBooksUsecase interface {
	GetByUid(uid string) (*models.Book, error)
	GetByLibraryUid(luid string, all bool) ([]*models.Book, error)
	UpdateBook(buid string, info *models.BookPatch) error
}

