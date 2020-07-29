package main

import (
	"backend/database"
	"backend/handlers"
	"backend/redis"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5431
	user     = "postgres"
	password = "docker"
	dbname   = "game_web_db"
)

func newRouter() *mux.Router {
	// Declare a router
	r := mux.NewRouter()

	// Declare static file directory
	staticFileDirectory := http.Dir("./static/")

	// Create static file server for our static files, i.e., .html, .css, etc
	staticFileServer := http.FileServer(staticFileDirectory)

	// Create file handler. Although the static files are placed inside `./static/` folder
	// in our local directory, it is served at the root (i.e., http://localhost:8080/)
	// when browsed in a browser. Hence, we need `http.StripPrefix` function to change the
	// file serve path.
	staticFileHandler := http.StripPrefix("/", staticFileServer)

	// Add static file handler to our router
	r.Handle("/", staticFileHandler).Methods("GET")

	// Add handler for `get` and `post` people functions
	r.HandleFunc("/user", handlers.GetUsersHandler).Methods("GET")
	r.HandleFunc("/user", handlers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/user", handlers.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/user", handlers.DeleteUserHandler).Methods("DELETE")

	r.HandleFunc("/game", handlers.GetGamesHandler).Methods("GET")
	r.HandleFunc("/game", handlers.СreateGameHandler).Methods("POST")
	r.HandleFunc("/game", handlers.UpdateGameHandler).Methods("PUT")
	r.HandleFunc("/game", handlers.DeleteGameHandler).Methods("DELETE")

	r.HandleFunc("/matches", handlers.GetMatchesHandler).Methods("GET")
	r.HandleFunc("/match", handlers.CreateMatchHandler).Methods("POST")
	r.HandleFunc("/match", handlers.UpdateMatchHandler).Methods("PUT")

	r.HandleFunc("/leaderboard", handlers.GetLeaderboardHandler).Methods("GET")
	r.HandleFunc("/leaderboard", handlers.CreateLeaderboardHandler).Methods("POST")
	r.HandleFunc("/leaderboard", handlers.UpdateLeaderboardHandler).Methods("PUT")
	r.HandleFunc("/leaderboard", handlers.DeleteLeaderboardHandler).Methods("DELETE")

	r.HandleFunc("/friendship", handlers.GetFriendshipHandler).Methods("GET")
	r.HandleFunc("/friendship", handlers.CreateFriendshipHandler).Methods("POST")
	r.HandleFunc("/friendship", handlers.DeleteFriendshipHandler).Methods("DELETE")

	r.HandleFunc("/move", handlers.GetMoveHandler).Methods("GET")
	r.HandleFunc("/move", handlers.CreateMoveHandler).Methods("POST")

	r.HandleFunc("/token", handlers.CreateTokenHandler).Methods("POST")
	r.HandleFunc("/token", handlers.CheckTokenHandler).Methods("GET")

	return r
}

func main() {

	var err error
	handlers.Store = database.InitStore(port, host, user, password, dbname)
	handlers.Redis, err = redis.InitRedis("tcp", "localhost:6380")

	if err != nil {
		log.Fatal(err.Error())
	}

	r := newRouter()

	log.Fatal(http.ListenAndServe(":8080", r))
}
