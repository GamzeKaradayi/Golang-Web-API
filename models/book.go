package models

type Book struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Writer   string `json:"writer"`
	Category string `json:"category"`
	InStock  bool   `json:"instock"`
}
