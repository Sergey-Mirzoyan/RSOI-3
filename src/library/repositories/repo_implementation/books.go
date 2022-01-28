package repo_implementation

import (
	"errors"
	"fmt"
	"lab2/src/database/pgdb"
	"lab2/src/library/models"
	"lab2/src/utility"
)

const (
	selectBooksByLibraryUid = `select id, book_uid, name, author, genre, condition from books where id in
								(select book_id from library_books where available_count > 0 and library_id in 
									(select id from library where library_uid=$1));`
	selectAllBooksByLibraryUid = `select id, book_uid, name, author, genre, condition from books where id in
								(select book_id from library_books where library_id in 
									(select id from library where library_uid=$1));`
	selectBookByUid = `select id, book_uid, name, author, genre, condition from books where book_uid=$1;`
	updateBookByUid = `update books set name=$2, author=$3, genre=$4, condition=$5 where book_uid=$1`
)

type BooksRepository struct {
	db *pgdb.DBManager
}

func NewBooksRepository(manager *pgdb.DBManager) *BooksRepository {
	return &BooksRepository{db: manager}
}

func (br *BooksRepository) GetByUid(uid string) (*models.Book, error) {
	data, err := br.db.Query(selectBookByUid, uid)
	if err != nil {
		fmt.Printf("Failed to get libraries from db\n")
	}

	res := &models.Book{
		Id:        utility.BytesToInt(data[0][0]),
		BookUid:   utility.BytesToUid(data[0][1]),
		Name:      utility.BytesToString(data[0][2]),
		Author:    utility.BytesToString(data[0][3]),
		Genre:     utility.BytesToString(data[0][4]),
		Condition: utility.BytesToString(data[0][5]),
	}
	return res, err
}

func (br *BooksRepository) GetByLibraryUid(luid string, all bool) ([]*models.Book, error) {
	var data [][][]byte
	var err error
	if all {
		data, err = br.db.Query(selectAllBooksByLibraryUid, luid)
	} else {
		data, err = br.db.Query(selectBooksByLibraryUid, luid)
	}
	if err != nil {
		fmt.Printf("Failed to get libraries from db\n")
	}
	var res []*models.Book

	for _, row := range data {
		res = append(res, &models.Book{
			Id:        utility.BytesToInt(row[0]),
			BookUid:   utility.BytesToUid(row[1]),
			Name:      utility.BytesToString(row[2]),
			Author:    utility.BytesToString(row[3]),
			Genre:     utility.BytesToString(row[4]),
			Condition: utility.BytesToString(row[5]),
		})
	}
	return res, err
}

func (br *BooksRepository) UpdateBook(buid string, info *models.BookPatch) error {
	affected, err := br.db.Exec(updateBookByUid, buid, info.Name, info.Author, info.Genre, info.Condition)
	if err != nil {
		fmt.Printf("Failed to update a book with uid=%s in db\n", buid)
		return err
	}
	if affected == 0 {
		err = errors.New("No book found in database that matches the received uid")
	}
	return err
}
