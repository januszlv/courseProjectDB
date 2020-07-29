package models

//Move ...
type Move struct {
	ID       string
	MatchID  string `json:"match_id"`
	Datetime string `json:"datetime"`
	Info     string `json:"info"`
}
