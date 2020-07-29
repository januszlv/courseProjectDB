package handlers

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//GetMoveHandler ...
func GetMoveHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve people from postgresql database using our `store` interface variable's
	// `func (*dbstore) GetUser` pointer receiver method defined in `store.go` file
	movesList, err := Store.GetMoves()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Convert the `personList` variable to JSON
	moveListBytes, err := json.Marshal(movesList)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Write JSON list of persons to response
	w.Write(moveListBytes)
}

//CreateMoveHandler ...
func CreateMoveHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
	defer r.Body.Close()

	move := models.Move{}
	err = json.Unmarshal(body, &move)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	id, datetime, err := Store.CreateMove(&move)
	move.ID = id
	move.Datetime = datetime

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}
