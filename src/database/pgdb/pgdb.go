package pgdb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DBManager struct {
	Pool *pgxpool.Pool
}

func NewPGDBManager() *DBManager {
	return &DBManager{Pool: nil}
}

func (db *DBManager) Connect(connectionString string) error {
	pool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return err
	}
	db.Pool = pool
	return nil
}

func (db *DBManager) Disconnect() {
	db.Pool.Close()
	db.Pool = nil
}

func (db *DBManager) Query(queryString string, params ...interface{}) ([][][]byte, error) {
	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, queryString, params...)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer rows.Close()

	result := make([][][]byte, 0)
	for rows.Next() {
		row := make([][]byte, 0)
		row = append(row, rows.RawValues()...)
		result = append(result, row)
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return result, nil
}

func (db *DBManager) Exec(queryString string, params ...interface{}) (int, error) {
	ctx := context.Background()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		fmt.Println("Error: ", err)
		return 0, err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, queryString, params...)
	if err != nil {
		fmt.Println("Error: ", err)
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Println("Error: ", err)
		return 0, err
	}

	return int(result.RowsAffected()), nil
}
