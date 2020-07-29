package models

//Leaderboard ...
type Leaderboard struct {
	ID     string
	UserID string `json:"user_id"`
	GameID string `json:"game_id"`
	Score  string `json:"score"`
}
