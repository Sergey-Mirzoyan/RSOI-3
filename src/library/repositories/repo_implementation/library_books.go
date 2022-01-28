package repo_implementation

import (
	"errors"
	"fmt"
	"lab2/src/database/pgdb"
	"lab2/src/utility"
)

const (
	updateLibraryBooks = `update library_books set available_count = $3 where library_id in
							(select id from library where library_uid=$1) and book_id in
							(select id from books where book_uid=$2);`
	selectLibraryBooks = `select available_count from library_books where library_id in
							(select id from library where library_uid=$1) and book_id in
							(select id from books where book_uid=$2);`
)

type LibraryBooksRepository struct {
	db *pgdb.DBManager
}

func NewLibraryBooksRepository(manager *pgdb.DBManager) *LibraryBooksRepository {
	return &LibraryBooksRepository{db: manager}
}

func (lbr *LibraryBooksRepository) UpdateBooksAmount(luid string, buid string, amount int) error {
	affected, err := lbr.db.Exec(updateLibraryBooks, luid, buid, amount)
	if err != nil {
		fmt.Printf("Failed to update library_books\n")
	}
	if affected != 1 {
		fmt.Printf("Affected %d rows while updating library_books (expected 1)\n", affected)
		err = errors.New("Bad library/book uid combination")
	}
	return err
}

func (lbr *LibraryBooksRepository) GetBooksAmount(luid string, buid string) (int, error) {
	data, err := lbr.db.Query(selectLibraryBooks, luid, buid)
	if err != nil {
		fmt.Printf("Failed to update library_books\n")
	}
	res := utility.BytesToInt(data[0][0])
	return res, err
}
