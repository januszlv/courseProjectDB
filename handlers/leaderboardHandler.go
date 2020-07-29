package handlers

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//GetLeaderboardHandler ...
func GetLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve people from postgresql database using our `store` interface variable's
	// `func (*dbstore) GetUser` pointer receiver method defined in `store.go` file
	leadsList, err := Store.GetLeaderboard()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Convert the `personList` variable to JSON
	leadsListBytes, err := json.Marshal(leadsList)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Write JSON list of persons to response
	w.Write(leadsListBytes)
}

//CreateLeaderboardHandler ...
func CreateLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
	defer r.Body.Close()

	leaderboard := models.Leaderboard{}
	err = json.Unmarshal(body, &leaderboard)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	id, err := Store.CreateLeaderboard(&leaderboard)
	leaderboard.ID = id

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}

//UpdateLeaderboardHandler ...
func UpdateLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println(fmt.Errorf("Value does not exist"))
		return
	}

	id := keys[0]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
	defer r.Body.Close()

	leaderboard := models.Leaderboard{}
	err = json.Unmarshal(body, &leaderboard)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	err = Store.UpdLeaderboard(&leaderboard, id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}

//DeleteLeaderboardHandler ...
func DeleteLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println(fmt.Errorf("Value does not exist"))
		return
	}

	id := keys[0]

	err := Store.DelLeaderboard(id)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}
