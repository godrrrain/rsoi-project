package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Library struct {
	ID          int    `json:"id"`
	Library_uid string `json:"library_uid"`
	Name        string `json:"name"`
	City        string `json:"city"`
	Address     string `json:"address"`
}

type Book struct {
	ID              int    `json:"id"`
	Book_uid        string `json:"book_uid"`
	Name            string `json:"name"`
	Author          string `json:"author"`
	Genre           string `json:"genre"`
	Condition       string `json:"condition"`
	Available_count int    `json:"available_count"`
}

type BookInfo struct {
	ID        int    `json:"id"`
	Book_uid  string `json:"book_uid"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	Genre     string `json:"genre"`
	Condition string `json:"condition"`
}

type Storage interface {
	// Insert(e *Person)
	// Get(id int) (Person, error)
	// Update(e *Person) error
	// Delete(id int) error
	// GetAll() []Person
	GetLibrariesByCity(ctx context.Context, city string) ([]Library, error)
	GetBooksByLibraryUid(ctx context.Context, libraryUid string, showAll bool) ([]Book, error)
	GetBookByUid(ctx context.Context, bookUid string) (Book, error)
	GetBookInfoByUid(ctx context.Context, bookUid string) (BookInfo, error)
	GetLibraryByUid(ctx context.Context, libraryUid string) (Library, error)
	UpdateBookCount(ctx context.Context, bookId int, count int) error
	UpdateBookCondition(ctx context.Context, bookUid string, condition string) error
}

type postgres struct {
	db *pgxpool.Pool
}

func NewPgStorage(ctx context.Context, connString string) (*postgres, error) {
	var pgInstance *postgres
	var pgOnce sync.Once
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			fmt.Printf("Unable to create connection pool: %v\n", err)
			return
		}

		pgInstance = &postgres{db}
	})

	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func (pg *postgres) GetLibrariesByCity(ctx context.Context, city string) ([]Library, error) {
	query := fmt.Sprintf(`SELECT id, library_uid, name, city, address FROM library WHERE city = '%s'`, city)

	rows, err := pg.db.Query(ctx, query)

	var libraries []Library

	if err != nil {
		return libraries, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	libraries, err = pgx.CollectRows(rows, pgx.RowToStructByName[Library])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return libraries, err
	}

	return libraries, nil
}

func (pg *postgres) GetBooksByLibraryUid(ctx context.Context, libraryUid string, showAll bool) ([]Book, error) {
	query := ""
	if showAll {
		query = fmt.Sprintf(`SELECT books.*, library_books.available_count from library_books, books, library 
	where library.library_uid = '%s' and library.id = library_books.library_id 
	and books.id = library_books.book_id;`, libraryUid)
	} else {
		query = fmt.Sprintf(`SELECT books.*, library_books.available_count from library_books, books, library 
	where library.library_uid = '%s' and library.id = library_books.library_id 
	and books.id = library_books.book_id and library_books.available_count > 0;`, libraryUid)
	}

	rows, err := pg.db.Query(ctx, query)

	var books []Book

	if err != nil {
		return books, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	books, err = pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return books, err
	}

	return books, nil
}

func (pg *postgres) GetBookByUid(ctx context.Context, bookUid string) (Book, error) {

	query := fmt.Sprintf(`SELECT books.*, library_books.available_count from library_books, books, library 
	where books.book_uid = '%s' and library.id = library_books.library_id 
	and books.id = library_books.book_id;`, bookUid)

	rows, err := pg.db.Query(ctx, query)

	var book Book

	if err != nil {
		return book, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	book, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[Book])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return book, err
	}

	return book, nil
}

func (pg *postgres) GetBookInfoByUid(ctx context.Context, bookUid string) (BookInfo, error) {

	query := fmt.Sprintf(`SELECT * FROM books WHERE book_uid = '%s'`, bookUid)

	rows, err := pg.db.Query(ctx, query)

	var book BookInfo

	if err != nil {
		return book, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	book, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[BookInfo])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return book, err
	}

	return book, nil
}

func (pg *postgres) GetLibraryByUid(ctx context.Context, libraryUid string) (Library, error) {

	query := fmt.Sprintf(`SELECT id, library_uid, name, city, address FROM library WHERE library_uid = '%s'`, libraryUid)

	rows, err := pg.db.Query(ctx, query)

	var library Library

	if err != nil {
		return library, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	library, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[Library])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return library, err
	}

	return library, nil
}

func (pg *postgres) UpdateBookCount(ctx context.Context, bookId int, count int) error {
	query := fmt.Sprintf(`UPDATE library_books SET available_count = %d WHERE book_id = %d`, count, bookId)

	_, err := pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func (pg *postgres) UpdateBookCondition(ctx context.Context, bookUid string, condition string) error {
	query := fmt.Sprintf(`UPDATE books SET condition = '%s' WHERE book_id = '%s'`, condition, bookUid)

	_, err := pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

// func (pg *postgres) UpdateBookCount(ctx context.Context, bookUid string) error {

// 	query := fmt.Sprintf(`SELECT id, library_uid, name, city, address FROM library WHERE library_uid = '%s'`, libraryUid)

// 	rows, err := pg.db.Query(ctx, query)

// 	var library Library

// 	if err != nil {
// 		fmt.Errorf("unable to query: %w", err)
// 	}
// 	defer rows.Close()

// 	library, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[Library])
// 	if err != nil {
// 		fmt.Printf("CollectRows error: %v", err)
// 		return err
// 	}

// 	query := fmt.Sprintf(`UPDATE library_books SET available_count = '%s' WHERE bookUid = '%s'`, condition, bookUid)

// 	_, err := pg.db.Exec(ctx, query)
// 	if err != nil {
// 		return fmt.Errorf("unable to insert row: %w", err)
// 	}

// 	return nil
// }
