package models

type Books struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Stock int `json:"stock"`
}