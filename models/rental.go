package models

import "time"

type Rentals struct {
	ID string `json:"id"`
	BookID string `json:"book_id"`
	UserID string `json:"user_id"`
	Username string `json:"username"`
	RentalDate time.Time `json:"rental_date"`
	ReturnDate time.Time `json:"return_date"`
	Status string `json:"status"`
	
}