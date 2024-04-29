package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Rating struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Stars    int    `json:"stars"`
}

type Storage interface {
	// Insert(e *Person)
	// Get(id int) (Person, error)
	// Update(e *Person) error
	// Delete(id int) error
	// GetAll() []Person
	GetRating(ctx context.Context, username string) (Rating, error)
	UpdateRating(ctx context.Context, username string, stars int) error
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

func (pg *postgres) GetRating(ctx context.Context, username string) (Rating, error) {
	query := fmt.Sprintf(`SELECT id, username, stars FROM rating WHERE username = '%s'`, username)

	rows, err := pg.db.Query(ctx, query)

	var rating Rating

	if err != nil {
		return rating, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	rating, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[Rating])

	if errors.Is(err, pgx.ErrNoRows) {
		return rating, errors.New("username not found")
	}

	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return rating, err
	}

	return rating, nil
}

func (pg *postgres) UpdateRating(ctx context.Context, username string, stars int) error {

	//temporary solution
	query := fmt.Sprintf(`SELECT id, username, stars FROM rating WHERE username = '%s'`, username)

	rows, err := pg.db.Query(ctx, query)

	var rating Rating

	if err != nil {
		return fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	rating, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[Rating])

	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("username not found")
	}

	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return err
	}

	resStars := rating.Stars + stars
	if resStars < 0 {
		resStars = 0
	}
	if resStars > 100 {
		resStars = 100
	}

	//upd
	query = fmt.Sprintf(`UPDATE rating SET stars = %d WHERE username = '%s'`, resStars, username)

	_, err = pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to update row: %w", err)
	}

	return nil
}
