package usecases

type IRatingUsecase interface {
	GetByUser(username string) (int, error)
	AlterByUser(diff int, username string) error
}

