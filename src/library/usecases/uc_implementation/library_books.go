package uc_implementation

import "lab2/src/library/repositories"

type LibraryBooksUsecase struct {
	lbr repositories.ILibraryBooksRepository
}

func NewLibraryBooksUsecase(repo repositories.ILibraryBooksRepository) *LibraryBooksUsecase {
	return &LibraryBooksUsecase{ lbr: repo }
}

func (lbc *LibraryBooksUsecase) UpdateBooksAmount(luid string, buid string, amount int) error {
	return lbc.lbr.UpdateBooksAmount(luid, buid, amount)
}

func (lbc *LibraryBooksUsecase) GetBooksAmount(luid string, buid string) (int, error) {
	return lbc.lbr.GetBooksAmount(luid, buid)
}