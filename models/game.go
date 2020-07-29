package models

//Game ...
type Game struct {
	ID          string
	Name        string `json:"name"`
	Description string `json:"description"`
}
