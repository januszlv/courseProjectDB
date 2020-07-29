package handlers

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//GetMatchesHandler ...
func GetMatchesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve people from postgresql database using our `store` interface variable's
	// `func (*dbstore) GetUser` pointer receiver method defined in `store.go` file
	matchList, err := Store.GetMatches()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Convert the `personList` variable to JSON
	matchListBytes, err := json.Marshal(matchList)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Write JSON list of persons to response
	w.Write(matchListBytes)
}

//CreateMatchHandler ...
func CreateMatchHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
	defer r.Body.Close()

	match := models.Match{}
	err = json.Unmarshal(body, &match)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	id, datetime, err := Store.CreateMatch(&match)
	match.ID = id
	match.Datetime = datetime

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}

//UpdateMatchHandler ...
func UpdateMatchHandler(w http.ResponseWriter, r *http.Request) {
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

	match := models.Match{}
	err = json.Unmarshal(body, &match)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	err = Store.UpdMatch(&match, id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}
