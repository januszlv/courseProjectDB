package models

//Friendship ...
type Friendship struct {
	ID      string
	UserID  string `json:"user_id"`
	User2ID string `json:"user2_id"`
}
