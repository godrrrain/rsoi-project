package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"lab2/src/library-service/storage"

	"github.com/gin-gonic/gin"
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
	Address     string `json:"address"`
	City        string `json:"city"`
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

type RequestUpdateReservation struct {
	Condition string `json:"condition"`
	Date      string `json:"date"`
}

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) GetLibrariesByCity(c *gin.Context) {

	libraries, err := h.storage.GetLibrariesByCity(context.Background(), c.Query("city"))

	if err != nil {
		fmt.Printf("failed to get libraries %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, LibrariesToResponse(libraries))
}

func (h *Handler) GetBooksByLibraryUid(c *gin.Context) {

	showAll, err := strconv.ParseBool(c.Query("showAll"))
	if err != nil {
		showAll = false
	}

	books, err := h.storage.GetBooksByLibraryUid(context.Background(), c.Param("uid"), showAll)

	if err != nil {
		fmt.Printf("failed to get libraries %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, BooksToResponse(books))
}

func (h *Handler) UpdateBookCount(c *gin.Context) {

	book, err := h.storage.GetBookByUid(context.Background(), c.Param("uid"))

	if err != nil {
		fmt.Printf("failed to get libraries %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	count := 1
	if c.Param("inc") == "1" {
		count = -1
	}

	err = h.storage.UpdateBookCount(context.Background(), book.ID, book.Available_count-count)

	if err != nil {
		fmt.Printf("failed to update book count %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "count updated",
	})
}

func (h *Handler) UpdateBookCondition(c *gin.Context) {

	book, err := h.storage.GetBookInfoByUid(context.Background(), c.Param("uid"))

	if err != nil {
		fmt.Printf("failed to get libraries %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var reqUpdRes RequestUpdateReservation

	err = json.NewDecoder(c.Request.Body).Decode(&reqUpdRes)
	if err != nil {
		fmt.Printf("failed to decode body %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if reqUpdRes.Condition != book.Condition {
		err = h.storage.UpdateBookCondition(context.Background(), c.Param("uid"), reqUpdRes.Condition)

		if err != nil {
			fmt.Printf("failed to update reservation %s\n", err.Error())
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, MessageResponse{
			Message: "condition updated",
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "condition already updated",
	})
}

func (h *Handler) GetBookInfoByUid(c *gin.Context) {

	book, err := h.storage.GetBookInfoByUid(context.Background(), c.Param("uid"))

	if err != nil {
		fmt.Printf("failed to get libraries %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, BookToUserResponse{
		Book_uid: book.Book_uid,
		Name:     book.Name,
		Author:   book.Author,
		Genre:    book.Genre,
	})
}

func (h *Handler) GetLibraryByUid(c *gin.Context) {

	library, err := h.storage.GetLibraryByUid(context.Background(), c.Param("uid"))

	if err != nil {
		fmt.Printf("failed to get libraries %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, LibraryToResponse(library))
}

func LibraryToResponse(library storage.Library) LibraryResponse {
	return LibraryResponse{
		Library_uid: library.Library_uid,
		Name:        library.Name,
		City:        library.City,
		Address:     library.Address,
	}
}

func LibrariesToResponse(libraries []storage.Library) []LibraryResponse {
	if libraries == nil {
		return nil
	}

	res := make([]LibraryResponse, len(libraries))

	for index, value := range libraries {
		res[index] = LibraryToResponse(value)
	}

	return res
}

func BookToResponse(book storage.Book) BookResponse {
	return BookResponse{
		Book_uid:        book.Book_uid,
		Name:            book.Name,
		Author:          book.Author,
		Genre:           book.Genre,
		Condition:       book.Condition,
		Available_count: book.Available_count,
	}
}

func BooksToResponse(books []storage.Book) []BookResponse {
	if books == nil {
		return nil
	}

	res := make([]BookResponse, len(books))

	for index, value := range books {
		res[index] = BookToResponse(value)
	}

	return res
}

func (h *Handler) GetHealth(c *gin.Context) {
	c.Status(http.StatusOK)
}
