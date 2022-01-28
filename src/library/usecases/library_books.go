package usecases

type ILibraryBooksUsecase interface {
	UpdateBooksAmount(luid string, buid string, amount int) error
	GetBooksAmount(luid string, buid string) (int, error)
}