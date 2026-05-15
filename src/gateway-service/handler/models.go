package handler

const (
	ratingService      string = "http://rating-service:8050"
	libraryService     string = "http://library-service:8060"
	reservationService string = "http://reservation-service:8070"
	statisticsService  string = "http://statistics-service:8040"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type LibraryResponse struct {
	Library_uid string `json:"libraryUid"`
	Name        string `json:"name"`
	City        string `json:"city"`
	Address     string `json:"address"`
}

type LibrariesLimited struct {
	Page          int               `json:"page"`
	PageSize      int               `json:"pageSize"`
	TotalElements int               `json:"totalElements"`
	Items         []LibraryResponse `json:"items"`
}

type BookResponse struct {
	Book_uid        string `json:"bookUid"`
	Name            string `json:"name"`
	Author          string `json:"author"`
	Genre           string `json:"genre"`
	Condition       string `json:"condition"`
	Available_count int    `json:"availableCount"`
}

type BookToUserResponse struct {
	Book_uid string `json:"bookUid"`
	Name     string `json:"name"`
	Author   string `json:"author"`
	Genre    string `json:"genre"`
}

type BookLimited struct {
	Page          int            `json:"page"`
	PageSize      int            `json:"pageSize"`
	TotalElements int            `json:"totalElements"`
	Items         []BookResponse `json:"items"`
}

type RatingResponse struct {
	Stars int `json:"stars"`
}

type UpdateRatingRequest struct {
	Stars    int    `json:"stars"`
	Username string `json:"username"`
}

type ReservationResponse struct {
	Reservation_uid string `json:"reservationUid"`
	Username        string `json:"username"`
	Book_uid        string `json:"bookUid"`
	Library_uid     string `json:"libraryUid"`
	Status          string `json:"status"`
	Start_date      string `json:"startDate"`
	Till_date       string `json:"tillDate"`
}

type ReservationToUserResponse struct {
	Reservation_uid string             `json:"reservationUid"`
	Status          string             `json:"status"`
	Start_date      string             `json:"startDate"`
	Till_date       string             `json:"tillDate"`
	Book            BookToUserResponse `json:"book"`
	Library         LibraryResponse    `json:"library"`
}

type TakeBookResponse struct {
	Reservation_uid string             `json:"reservationUid"`
	Status          string             `json:"status"`
	Start_date      string             `json:"startDate"`
	Till_date       string             `json:"tillDate"`
	Book            BookToUserResponse `json:"book"`
	Library         LibraryResponse    `json:"library"`
	Rating          RatingResponse     `json:"rating"`
}

type CreateReservationRequest struct {
	BookUid    string `json:"bookUid"`
	LibraryUid string `json:"libraryUid"`
	TillDate   string `json:"tillDate"`
}

type UpdateReservationRequest struct {
	Condition string `json:"condition"`
	Date      string `json:"date"`
}

type ReservationAmount struct {
	Amount int `json:"amount"`
}

type ReservationsResponse struct {
	Reservation_uid string             `json:"reservationUid"`
	Status          string             `json:"status"`
	Start_date      string             `json:"startDate"`
	Till_date       string             `json:"tillDate"`
	Book            BookToUserResponse `json:"book"`
	Library         LibraryResponse    `json:"library"`
	Username        Username           `json:"username"`
}

type Username struct {
	Name string `json:"name"`
}
