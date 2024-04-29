package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Reservation struct {
	ID              int       `json:"id"`
	Reservation_uid string    `json:"reservation_uid"`
	Username        string    `json:"username"`
	Book_uid        string    `json:"book_uid"`
	Library_uid     string    `json:"library_uid"`
	Status          string    `json:"status"`
	Start_date      time.Time `json:"start_date"`
	Till_date       time.Time `json:"till_date"`
}

type ReservationAmount struct {
	Amount int `json:"amount"`
}

type Storage interface {
	// Update(e *Person) error
	// Delete(id int) error
	// GetAll() []Person
	GetReservations(ctx context.Context, username string) ([]Reservation, error)
	GetReservationByUid(ctx context.Context, reservation_uid string) (Reservation, error)
	GetRentedReservationAmount(ctx context.Context, username string) (ReservationAmount, error)
	CreateReservation(ctx context.Context, username string, bookUid string, libraryUid string, tillDate string) (Reservation, error)
	UpdateReservationStatus(ctx context.Context, reservation_uid string, status string) error
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

func (pg *postgres) CreateReservation(ctx context.Context, username string, bookUid string, libraryUid string, tillDate string) (Reservation, error) {

	var reservation Reservation

	uid := uuid.New()

	reservation_uid := uid.String()

	start_date := time.Now().UTC().Format("2006-01-02")

	query := `INSERT INTO reservation (reservation_uid, username, book_uid, library_uid, status, start_date, till_date) 
	VALUES (@reservation_uid, @username, @book_uid, @library_uid, @status, @start_date, @till_date)`
	args := pgx.NamedArgs{
		"reservation_uid": reservation_uid,
		"username":        username,
		"book_uid":        bookUid,
		"library_uid":     libraryUid,
		"status":          "RENTED",
		"start_date":      start_date,
		"till_date":       tillDate,
	}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return reservation, fmt.Errorf("unable to insert row: %w", err)
	}

	tillDateTime, err := time.Parse("2006-01-02", tillDate)
	if err != nil {
		fmt.Println(err)
		return reservation, fmt.Errorf("unable to convert time: %w", err)
	}

	reservation.Reservation_uid = reservation_uid
	reservation.Username = username
	reservation.Book_uid = bookUid
	reservation.Library_uid = libraryUid
	reservation.Status = "RENTED"
	reservation.Start_date = time.Now().UTC()
	reservation.Till_date = tillDateTime

	return reservation, nil
}

func (pg *postgres) GetReservationByUid(ctx context.Context, reservation_uid string) (Reservation, error) {

	query := fmt.Sprintf(`SELECT * FROM reservation WHERE reservation_uid = '%s'`, reservation_uid)

	rows, err := pg.db.Query(ctx, query)

	var reservation Reservation

	if err != nil {
		return reservation, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	reservation, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return reservation, err
	}

	return reservation, nil
}

func (pg *postgres) GetReservations(ctx context.Context, username string) ([]Reservation, error) {

	query := fmt.Sprintf(`SELECT * FROM reservation WHERE username = '%s'`, username)

	rows, err := pg.db.Query(ctx, query)

	var reservations []Reservation

	if err != nil {
		return reservations, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	reservations, err = pgx.CollectRows(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return reservations, err
	}

	return reservations, nil
}

func (pg *postgres) GetRentedReservationAmount(ctx context.Context, username string) (ReservationAmount, error) {

	query := fmt.Sprintf(`SELECT * FROM reservation WHERE username = '%s' and status = 'RENTED'`, username)

	rows, err := pg.db.Query(ctx, query)

	var reservationAmount ReservationAmount

	if err != nil {
		return reservationAmount, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	reservations, err := pgx.CollectRows(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return reservationAmount, err
	}
	reservationAmount.Amount = len(reservations)

	return reservationAmount, nil
}

func (pg *postgres) UpdateReservationStatus(ctx context.Context, reservation_uid string, status string) error {
	query := fmt.Sprintf(`UPDATE reservation SET status = '%s' WHERE reservation_uid = '%s'`, status, reservation_uid)

	_, err := pg.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}
