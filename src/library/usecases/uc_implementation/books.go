package uc_implementation

import (
"lab2/src/library/models"
"lab2/src/library/repositories"
)

type BooksUsecase struct {
	br repositories.IBooksRepository
}

func NewBooksUsecase(repo repositories.IBooksRepository) *BooksUsecase {
	return &BooksUsecase{br: repo}
}

func (bc *BooksUsecase) GetByUid(uid string) (*models.Book, error) {
	return bc.br.GetByUid(uid)
}

func (bc *BooksUsecase) GetByLibraryUid(luid string, all bool) ([]*models.Book, error) {
	return bc.br.GetByLibraryUid(luid, all)
}

func (bc *BooksUsecase) UpdateBook(buid string, info *models.BookPatch) error {
	oldInfo, err := bc.br.GetByUid(buid)
	if err != nil {
		return err
	}
	if info.Name == nil {
		info.Name = &oldInfo.Name
	}
	if info.Author == nil {
		info.Author = &oldInfo.Author
	}
	if info.Genre == nil {
		info.Genre = &oldInfo.Genre
	}
	if info.Condition == nil {
		info.Condition = &oldInfo.Condition
	}
	return bc.br.UpdateBook(buid, info)
}
