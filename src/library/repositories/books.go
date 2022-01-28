package repositories

import "lab2/src/library/models"

type IBooksRepository interface {
	GetByUid(uid string) (*models.Book, error)
	GetByLibraryUid(luid string, all bool) ([]*models.Book, error)
	UpdateBook(buid string, info *models.BookPatch) error
}
