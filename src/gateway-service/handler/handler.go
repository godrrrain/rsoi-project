package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lab2/src/jobqueue"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
)

const (
	ratingService      string = "http://rating-service:8050"
	libraryService     string = "http://library-service:8060"
	reservationService string = "http://reservation-service:8070"
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

type Handler struct {
	libraryCB     *gobreaker.CircuitBreaker
	ratingCB      *gobreaker.CircuitBreaker
	reservationCB *gobreaker.CircuitBreaker
	jobScheduler  *jobqueue.JobScheduler
}

func NewHandler(libraryCircuitBreaker, ratingCircuitBreaker, reservationCircuitBreaker *gobreaker.CircuitBreaker, jobScheduler *jobqueue.JobScheduler) *Handler {
	return &Handler{
		libraryCB:     libraryCircuitBreaker,
		ratingCB:      ratingCircuitBreaker,
		reservationCB: reservationCircuitBreaker,
		jobScheduler:  jobScheduler,
	}
}

func (h *Handler) GetLibrariesByCity(c *gin.Context) {
	params := c.Request.URL.Query()
	requestURL := fmt.Sprintf("%s/api/v1/libraries/", libraryService)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	q := req.URL.Query()
	q.Add("city", c.Query("city"))
	req.URL.RawQuery = q.Encode()

	ires, err := h.libraryCB.Execute(func() (interface{}, error) {
		return http.DefaultClient.Do(req)
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Message: "Library Service unavailable"})
		return
	}

	res, ok := ires.(*http.Response)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var libraries []LibraryResponse
	if json.Unmarshal(resBody, &libraries) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	fmt.Println(libraries)

	pageParam := params.Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	sizeParam := params.Get("size")
	if sizeParam == "" {
		sizeParam = "100"
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	right := page * size
	if len(libraries) < right {
		right = len(libraries)
	}

	librariesStripped := make([]LibraryResponse, 0)

	if (page-1)*size <= len(libraries) {
		librariesStripped = libraries[(page-1)*size : right]
	}

	data := LibrariesLimited{
		Page:          page,
		PageSize:      size,
		TotalElements: len(librariesStripped),
		Items:         librariesStripped,
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetBooksByLibraryUid(c *gin.Context) {
	params := c.Request.URL.Query()
	requestURL := fmt.Sprintf("%s/api/v1/libraries/%s/books/", libraryService, c.Param("uid"))

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	q := req.URL.Query()
	q.Add("showAll", c.Query("showAll"))
	req.URL.RawQuery = q.Encode()

	ires, err := h.libraryCB.Execute(func() (interface{}, error) {
		return http.DefaultClient.Do(req)
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Message: "Library Service unavailable"})
		return
	}

	res, ok := ires.(*http.Response)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var books []BookResponse
	if json.Unmarshal(resBody, &books) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	fmt.Println(books)

	pageParam := params.Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	sizeParam := params.Get("size")
	if sizeParam == "" {
		sizeParam = "100"
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	right := page * size
	if len(books) < right {
		right = len(books)
	}

	booksStripped := make([]BookResponse, 0)

	if (page-1)*size <= len(books) {
		booksStripped = books[(page-1)*size : right]
	}

	data := BookLimited{
		Page:          page,
		PageSize:      size,
		TotalElements: len(booksStripped),
		Items:         booksStripped,
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetRating(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	requestURL := fmt.Sprintf("%s/api/v1/rating/", ratingService)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	req.Header.Set("X-User-Name", username)

	ires, err := h.ratingCB.Execute(func() (interface{}, error) {
		return http.DefaultClient.Do(req)
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Message: "Bonus Service unavailable"})
		return
	}

	res, ok := ires.(*http.Response)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var rating RatingResponse
	if json.Unmarshal(resBody, &rating) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RatingResponse{
		Stars: rating.Stars,
	})
}

func (h *Handler) GetReservations(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	requestURL := fmt.Sprintf("%s/api/v1/reservations/", reservationService)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	req.Header.Set("X-User-Name", username)

	ires, err := h.reservationCB.Execute(func() (interface{}, error) {
		return http.DefaultClient.Do(req)
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Message: "Reservation Service unavailable"})
		return
	}

	res, ok := ires.(*http.Response)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var reservations []ReservationResponse
	if json.Unmarshal(resBody, &reservations) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	response := make([]ReservationToUserResponse, len(reservations))

	for i, reservation := range reservations {
		requestBookURL := fmt.Sprintf("%s/api/v1/books/%s/", libraryService, reservation.Book_uid)

		req, err := http.NewRequest(http.MethodGet, requestBookURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		ires, err := h.libraryCB.Execute(func() (interface{}, error) {
			return http.DefaultClient.Do(req)
		})
		if err != nil {

			var book BookToUserResponse
			book.Book_uid = reservation.Book_uid

			var library LibraryResponse
			library.Library_uid = reservation.Library_uid

			response[i] = ReservationToUserResponse{
				Reservation_uid: reservation.Reservation_uid,
				Status:          reservation.Status,
				Start_date:      reservation.Start_date,
				Till_date:       reservation.Till_date,
				Book:            book,
				Library:         library,
			}

			fmt.Println("Library Service unavailable")
			continue
		}

		res, ok := ires.(*http.Response)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		var book BookToUserResponse
		if json.Unmarshal(resBody, &book) != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		requestLibraryURL := fmt.Sprintf("%s/api/v1/libraries/%s/", libraryService, reservation.Library_uid)

		reqLib, err := http.NewRequest(http.MethodGet, requestLibraryURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		iresLib, err := h.libraryCB.Execute(func() (interface{}, error) {
			return http.DefaultClient.Do(reqLib)
		})
		if err != nil {
			var book BookToUserResponse
			book.Book_uid = reservation.Book_uid

			var library LibraryResponse
			library.Library_uid = reservation.Library_uid

			response[i] = ReservationToUserResponse{
				Reservation_uid: reservation.Reservation_uid,
				Status:          reservation.Status,
				Start_date:      reservation.Start_date,
				Till_date:       reservation.Till_date,
				Book:            book,
				Library:         library,
			}

			fmt.Println("Library Service unavailable")
			continue
		}

		resLib, ok := iresLib.(*http.Response)
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		resLibBody, err := io.ReadAll(resLib.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		var library LibraryResponse
		if json.Unmarshal(resLibBody, &library) != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		response[i] = ReservationToUserResponse{
			Reservation_uid: reservation.Reservation_uid,
			Status:          reservation.Status,
			Start_date:      reservation.Start_date,
			Till_date:       reservation.Till_date,
			Book:            book,
			Library:         library,
		}

	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateReservation(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	var inputCreateBody CreateReservationRequest

	err := json.NewDecoder(c.Request.Body).Decode(&inputCreateBody)
	if err != nil {
		fmt.Printf("failed to decode body %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	//getting amount
	requestAmountURL := fmt.Sprintf("%s/api/v1/reservations/amount", reservationService)

	reqAmount, err := http.NewRequest(http.MethodGet, requestAmountURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	reqAmount.Header.Set("X-User-Name", username)

	iresAmount, err := h.reservationCB.Execute(func() (interface{}, error) {
		return http.DefaultClient.Do(reqAmount)
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Message: "Reservation Service unavailable"})
		return
	}

	resAmount, ok := iresAmount.(*http.Response)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resBodyAmount, err := io.ReadAll(resAmount.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var reservationAmount ReservationAmount
	if json.Unmarshal(resBodyAmount, &reservationAmount) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	//getting a rating
	requestRatingURL := fmt.Sprintf("%s/api/v1/rating/", ratingService)

	reqRating, err := http.NewRequest(http.MethodGet, requestRatingURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	reqRating.Header.Set("X-User-Name", username)

	iresRating, err := h.ratingCB.Execute(func() (interface{}, error) {
		return http.DefaultClient.Do(reqRating)
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Message: "Bonus Service unavailable"})
		return
	}

	resRating, ok := iresRating.(*http.Response)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resBodyRating, err := io.ReadAll(resRating.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var rating RatingResponse
	if json.Unmarshal(resBodyRating, &rating) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if reservationAmount.Amount >= rating.Stars {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "user cannot take new book",
		})
		return
	}

	//create reservation
	requestCreateURL := fmt.Sprintf("%s/api/v1/reservations", reservationService)

	marshalled, err := json.Marshal(inputCreateBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}

	reqCreate, err := http.NewRequest(http.MethodPost, requestCreateURL, bytes.NewReader(marshalled))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	reqCreate.Header.Set("X-User-Name", username)

	resCreate, err := http.DefaultClient.Do(reqCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resBodyCreate, err := io.ReadAll(resCreate.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var createReserv ReservationResponse
	if json.Unmarshal(resBodyCreate, &createReserv) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	//create response
	requestBookURL := fmt.Sprintf("%s/api/v1/books/%s/", libraryService, createReserv.Book_uid)

	reqBook, err := http.NewRequest(http.MethodGet, requestBookURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	iresBook, err := h.libraryCB.Execute(func() (interface{}, error) {
		return http.DefaultClient.Do(reqBook)
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Message: "Library Service unavailable"})
		return
	}

	resBook, ok := iresBook.(*http.Response)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resBodyBook, err := io.ReadAll(resBook.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var book BookToUserResponse
	if json.Unmarshal(resBodyBook, &book) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	requestLibraryURL := fmt.Sprintf("%s/api/v1/libraries/%s/", libraryService, createReserv.Library_uid)

	reqLib, err := http.NewRequest(http.MethodGet, requestLibraryURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	iresLib, err := h.libraryCB.Execute(func() (interface{}, error) {
		return http.DefaultClient.Do(reqLib)
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Message: "Library Service unavailable"})
		return
	}

	resLib, ok := iresLib.(*http.Response)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resLibBody, err := io.ReadAll(resLib.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var library LibraryResponse
	if json.Unmarshal(resLibBody, &library) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	response := TakeBookResponse{
		Reservation_uid: createReserv.Reservation_uid,
		Status:          createReserv.Status,
		Start_date:      createReserv.Start_date,
		Till_date:       createReserv.Till_date,
		Book:            book,
		Library:         library,
		Rating:          rating,
	}

	//update count
	requestUpdateCountURL := fmt.Sprintf("%s/api/v1/books/%s/count/0", libraryService, book.Book_uid)

	reqCount, err := http.NewRequest(http.MethodPut, requestUpdateCountURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resCount, err := http.DefaultClient.Do(reqCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if resCount.StatusCode != 200 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "error while updating count",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) ReturnBook(c *gin.Context) {

	resFee := 0

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	var inputUpdateBody UpdateReservationRequest

	err := json.NewDecoder(c.Request.Body).Decode(&inputUpdateBody)
	if err != nil {
		fmt.Printf("failed to decode body %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	//getting reservation info
	requestReservURL := fmt.Sprintf("%s/api/v1/reservations/info/%s", reservationService, c.Param("uid"))

	reqReserv, err := http.NewRequest(http.MethodGet, requestReservURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	reqReserv.Header.Set("X-User-Name", username)

	iresReserv, err := h.reservationCB.Execute(func() (interface{}, error) {
		return http.DefaultClient.Do(reqReserv)
	})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{Message: "Reservation Service unavailable"})
		return
	}

	resReserv, ok := iresReserv.(*http.Response)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resBodyReserv, err := io.ReadAll(resReserv.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var reservation ReservationResponse
	if json.Unmarshal(resBodyReserv, &reservation) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	//updating status
	requestStatusURL := fmt.Sprintf("%s/api/v1/reservations/%s", reservationService, c.Param("uid"))

	marshalled, err := json.Marshal(inputUpdateBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}

	reqStatus, err := http.NewRequest(http.MethodPut, requestStatusURL, bytes.NewReader(marshalled))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resStatus, err := http.DefaultClient.Do(reqStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if resStatus.StatusCode == 204 {
		resFee = resFee + 1
	}

	//updating condition
	requestConditionURL := fmt.Sprintf("%s/api/v1/books/%s/condition", libraryService, reservation.Book_uid)

	reqCondition, err := http.NewRequest(http.MethodPut, requestConditionURL, bytes.NewReader(marshalled))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resCondition, err := http.DefaultClient.Do(reqCondition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if resCondition.StatusCode == 201 {
		resFee = resFee + 1
	}

	//updating count
	requestCountURL := fmt.Sprintf("%s/api/v1/books/%s/count/1/", libraryService, reservation.Book_uid)

	reqCount, err := http.NewRequest(http.MethodPut, requestCountURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	_, err = http.DefaultClient.Do(reqCount)
	if err != nil {
		job := jobqueue.NewExecJob(func() (interface{}, error) {
			return http.DefaultClient.Do(reqCount)
		})
		h.jobScheduler.JobQueue <- job
	}

	//update rating
	var rating RatingResponse

	if resFee != 0 {
		resFee = resFee * -10
	} else {
		resFee = 1
	}

	rating.Stars = resFee

	requestUpdRatingURL := fmt.Sprintf("%s/api/v1/rating/", ratingService)

	marshalled, err = json.Marshal(rating)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}

	reqUpdRating, err := http.NewRequest(http.MethodPut, requestUpdRatingURL, bytes.NewReader(marshalled))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	reqUpdRating.Header.Set("X-User-Name", username)

	_, err = http.DefaultClient.Do(reqUpdRating)
	if err != nil {
		c.Status(http.StatusNoContent)
		job := jobqueue.NewExecJob(func() (interface{}, error) {
			return http.DefaultClient.Do(reqUpdRating)
		})
		h.jobScheduler.JobQueue <- job
		return
	}

	c.JSON(http.StatusNoContent, MessageResponse{
		Message: "Book was successfully returned",
	})
}

func (h *Handler) GetHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}
