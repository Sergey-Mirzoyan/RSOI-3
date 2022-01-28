package repositories

type IRatingRepository interface {
	GetByUser(username string) (int, error)
	SetByUser(value int, username string) error
}
