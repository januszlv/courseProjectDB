package models

//Match ...
type Match struct {
	ID       string
	UserID   string `json:"user_id"`
	GameID   string `json:"game_id"`
	Score    string `json:"score"`
	Datetime string `json:"datetime"`
}
