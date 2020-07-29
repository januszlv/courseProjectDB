package database

import (
	"backend/models"
	"database/sql"
	"fmt"
	"strconv"
)

/*Store is an interface for storing userdata;
* CreateUser -- add a new user;
* GetUsers --  get all existing users;
 */
type Store interface {
	CreateUser(user *models.User) (string, error)
	GetUsers(id string) ([]*models.User, error)
	DelUser(id string) error
	UpdUser(user *models.User, id string) error

	CreateGame(game *models.Game) (string, error)
	GetGames(id string) ([]*models.Game, error)
	DelGame(id string) error
	UpdGame(game *models.Game, id string) error

	CreateMatch(match *models.Match) (string, string, error)
	GetMatches() ([]*models.Match, error)
	UpdMatch(match *models.Match, id string) error

	CreateLeaderboard(leaderboard *models.Leaderboard) (string, error)
	GetLeaderboard() ([]*models.Leaderboard, error)
	DelLeaderboard(id string) error
	UpdLeaderboard(leaderboard *models.Leaderboard, id string) error

	CreateFriendship(friendship *models.Friendship) (string, error)
	GetFriendship(id string) ([]*models.Friendship, error)
	DelFriendship(id string) error

	CreateMove(match *models.Move) (string, string, error)
	GetMoves() ([]*models.Move, error)
}

// `dbStore` struct implements the `Store` interface. Variable `db` takes the pointer
// to the SQL database connection object.
type dbStore struct {
	db *sql.DB
}

//InitStore ...
func InitStore(port int, host, user, password, dbname string) Store {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	return &dbStore{db: db}
}

func (store *dbStore) CreateUser(user *models.User) (string, error) {
	id := 0
	err := store.db.QueryRow(
		"INSERT INTO \"user\" (user_name, full_name, password, motto) VALUES ($1, $2, $3, $4) RETURNING id",
		user.UserName, user.FullName, user.Password, user.Motto).Scan(&id)
	fmt.Println("User inserted")
	return string(id), err
}

func (store *dbStore) GetUsers(id string) ([]*models.User, error) {
	var rows *sql.Rows
	var err error

	if id != "" {
		rows, err = store.db.Query("SELECT id, user_name, full_name, password, motto FROM \"user\" WHERE id = $1 AND deleted IS NOT TRUE; ", id)
	} else {
		rows, err = store.db.Query("SELECT id, user_name, full_name, password, motto FROM \"user\" WHERE deleted IS NOT TRUE; ")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userList := []*models.User{}

	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.UserName, &user.FullName, &user.Password, &user.Motto); err != nil {
			return nil, err
		}

		userList = append(userList, user)
	}
	return userList, nil
}

func (store *dbStore) DelUser(id string) error {
	if id == "" {
		fmt.Println("ID is empty")
	}
	_, err := store.db.Query(
		"DELETE FROM \"user\" WHERE id = $1; ",
		id)
	fmt.Println("User deleted")
	return err
}

func (store *dbStore) UpdUser(user *models.User, id string) error {
	query := "UPDATE \"user\" SET "
	args := []interface{}{}
	if user.UserName != "" {
		query += "user_name = $1"
		args = append(args, user.UserName)
	}

	if user.FullName != "" {
		if len(args) == 0 {
			query += "full_name = $1"
		} else {
			count := len(args) + 1
			query += ", full_name = $" + string(count)
		}
		args = append(args, user.FullName)

	}

	if user.Password != "" {
		if len(args) == 0 {
			query += "password = $1"
		} else {
			count := len(args) + 1
			query += ", password = $" + string(count)
		}
		args = append(args, user.Password)

	}

	if user.Motto != "" {
		if len(args) == 0 {
			query += "motto = $1"
		} else {
			count := len(args) + 1
			query += ", motto = $" + string(count)
		}
		args = append(args, user.Motto)

	}

	query += "WHERE id = " + id

	_, err := store.db.Query(query, args...)
	fmt.Println("User updated")

	return err
}

func (store *dbStore) CreateGame(game *models.Game) (string, error) {
	id := 0
	err := store.db.QueryRow(
		"INSERT INTO game (name, description) VALUES ($1, $2) RETURNING id ",
		game.Name, game.Description).Scan(&id)
	fmt.Println("Game inserted")
	return string(id), err
}

func (store *dbStore) GetGames(id string) ([]*models.Game, error) {
	var rows *sql.Rows
	var err error

	if id != "" {
		rows, err = store.db.Query("SELECT id, name, description FROM game WHERE id = $1 AND deleted IS NOT TRUE; ", id)
	} else {
		rows, err = store.db.Query("SELECT id, name, description FROM game WHERE deleted IS NOT TRUE; ")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gameList := []*models.Game{}

	for rows.Next() {
		game := &models.Game{}
		if err := rows.Scan(&game.ID, &game.Name, &game.Description); err != nil {
			return nil, err
		}

		gameList = append(gameList, game)
	}
	return gameList, nil
}

func (store *dbStore) DelGame(id string) error {
	if id == "" {
		fmt.Println("ID is empty")
	}
	_, err := store.db.Query(
		"DELETE FROM game WHERE id = $1; ",
		id)
	fmt.Println("Game deleted")
	return err
}

func (store *dbStore) UpdGame(game *models.Game, id string) error {
	query := "UPDATE game SET "
	args := []interface{}{}
	if game.Name != "" {
		query += "user_name = $1"
		args = append(args, game.Name)
	}

	if game.Description != "" {
		if len(args) == 0 {
			query += "description = $1"
		} else {
			query += ", description = $2"
		}
		args = append(args, game.Description)

	}

	query += "WHERE id = " + id

	_, err := store.db.Query(query, args...)
	fmt.Println("Game updated")

	return err
}

func (store *dbStore) CreateMatch(match *models.Match) (string, string, error) {
	id := 0
	datetime := ""

	userID, err := strconv.Atoi(match.UserID)
	if err != nil {
		fmt.Println(fmt.Errorf("Atoi userID error: %v", err))

	}
	gameID, err := strconv.Atoi(match.GameID)
	if err != nil {
		fmt.Println(fmt.Errorf("Atoi gameID error: %v", err))

	}

	score, err := strconv.Atoi(match.Score)
	if err != nil {
		fmt.Println(fmt.Errorf("Atoi score error: %v", err))

	}
	err = store.db.QueryRow(
		"INSERT INTO match (user_id, game_id, score) VALUES ($1, $2, $3) RETURNING id, datetime",
		userID, gameID, score).Scan(&id, &datetime)
	fmt.Println("DATETIME IS: " + datetime)
	return string(id), datetime, err
}

func (store *dbStore) GetMatches() ([]*models.Match, error) {
	rows, err := store.db.Query("SELECT id, user_id, game_id, score, datetime FROM match; ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matchList := []*models.Match{}

	for rows.Next() {
		match := &models.Match{}
		if err := rows.Scan(&match.ID, &match.UserID, &match.GameID, &match.Score, &match.Datetime); err != nil {
			return nil, err
		}

		matchList = append(matchList, match)
	}
	return matchList, nil
}

func (store *dbStore) UpdMatch(match *models.Match, id string) error {
	query := "UPDATE match SET "

	if match.ID != "" || match.UserID != "" || match.GameID != "" {
		fmt.Println("ID updating is prohibited")
		return fmt.Errorf("ID updating is prohibited")
	}

	if string(match.Score) != "" {
		query += "score = $1 WHERE id = " + id
	}

	_, err := store.db.Query(query, match.Score)
	fmt.Println("Match updated")

	return err
}

func (store *dbStore) CreateLeaderboard(leaderboard *models.Leaderboard) (string, error) {
	id := 0

	userID, err := strconv.Atoi(leaderboard.UserID)
	if err != nil {
		fmt.Println(fmt.Errorf("Atoi userID error: %v", err))

	}
	gameID, err := strconv.Atoi(leaderboard.GameID)
	if err != nil {
		fmt.Println(fmt.Errorf("Atoi gameID error: %v", err))

	}

	score, err := strconv.Atoi(leaderboard.Score)
	if err != nil {
		fmt.Println(fmt.Errorf("Atoi score error: %v", err))

	}
	err = store.db.QueryRow(
		"INSERT INTO leaderboard (user_id, game_id, score) VALUES ($1, $2, $3) RETURNING id",
		userID, gameID, score).Scan(&id)
	return string(id), err
}

func (store *dbStore) GetLeaderboard() ([]*models.Leaderboard, error) {
	rows, err := store.db.Query("SELECT id, user_id, game_id, score FROM leaderboard; ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leadsList := []*models.Leaderboard{}

	for rows.Next() {
		leader := &models.Leaderboard{}
		if err := rows.Scan(&leader.ID, &leader.UserID, &leader.GameID, &leader.Score); err != nil {
			return nil, err
		}

		leadsList = append(leadsList, leader)
	}
	return leadsList, nil
}

func (store *dbStore) DelLeaderboard(id string) error {
	if id == "" {
		fmt.Println("ID is empty")
	}
	_, err := store.db.Query(
		"DELETE FROM leaderboard WHERE id = $1; ",
		id)
	fmt.Println("Leader deleted")
	return err
}
func (store *dbStore) UpdLeaderboard(leaderboard *models.Leaderboard, id string) error {
	query := "UPDATE leaderboard SET "

	if leaderboard.ID != "" || leaderboard.UserID != "" || leaderboard.GameID != "" {
		fmt.Println("ID updating is prohibited")
		return fmt.Errorf("ID updating is prohibited")
	}

	if string(leaderboard.Score) != "" {
		query += "score = $1 WHERE id = " + id
	}

	_, err := store.db.Query(query, leaderboard.Score)
	fmt.Println("Leaderboard updated")

	return err
}

func (store *dbStore) CreateFriendship(friendship *models.Friendship) (string, error) {
	id := 0

	userID, err := strconv.Atoi(friendship.UserID)
	if err != nil {
		fmt.Println(fmt.Errorf("Atoi userID error: %v", err))

	}
	user2ID, err := strconv.Atoi(friendship.User2ID)
	if err != nil {
		fmt.Println(fmt.Errorf("Atoi user2ID error: %v", err))

	}

	err = store.db.QueryRow(
		"INSERT INTO friendship (user_id, user2_id) VALUES ($1, $2) RETURNING id",
		userID, user2ID).Scan(&id)
	return string(id), err
}

func (store *dbStore) GetFriendship(id string) ([]*models.Friendship, error) {
	var rows *sql.Rows
	var err error

	if id != "" {
		rows, err = store.db.Query("SELECT id, user_id, user2_id FROM friendship WHERE id = $1", id)
	} else {
		rows, err = store.db.Query("SELECT id, user_id, user2_id FROM friendship; ")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friendList := []*models.Friendship{}

	for rows.Next() {
		friendship := &models.Friendship{}
		if err := rows.Scan(&friendship.ID, &friendship.UserID, &friendship.User2ID); err != nil {
			return nil, err
		}

		friendList = append(friendList, friendship)
	}
	return friendList, nil
}
func (store *dbStore) DelFriendship(id string) error {
	if id == "" {
		fmt.Println("ID is empty")
	}
	_, err := store.db.Query(
		"DELETE FROM friendship WHERE id = $1; ",
		id)
	fmt.Println("Friend deleted")
	return err
}

func (store *dbStore) CreateMove(move *models.Move) (string, string, error) {
	id := 0
	datetime := ""

	matchID, err := strconv.Atoi(move.MatchID)
	if err != nil {
		fmt.Println(fmt.Errorf("Atoi matchID error: %v", err))

	}

	err = store.db.QueryRow(
		"INSERT INTO move (match_id, info) VALUES ($1, $2) RETURNING id, datetime",
		matchID, move.Info).Scan(&id, &datetime)
	fmt.Println("DATETIME IS: " + datetime)
	return string(id), datetime, err
}

func (store *dbStore) GetMoves() ([]*models.Move, error) {
	rows, err := store.db.Query("SELECT id, match_id, datetime, info FROM move; ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	moveList := []*models.Move{}

	for rows.Next() {
		move := &models.Move{}
		if err := rows.Scan(&move.ID, &move.MatchID, &move.Datetime, &move.Info); err != nil {
			return nil, err
		}

		moveList = append(moveList, move)
	}
	return moveList, nil
}
