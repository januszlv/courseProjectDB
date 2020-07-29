package handlers

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//GetFriendshipHandler ...
func GetFriendshipHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve people from postgresql database using our `store` interface variable's
	// `func (*dbstore) GetUser` pointer receiver method defined in `store.go` file
	id := ""
	keys, ok := r.URL.Query()["id"]
	if !ok {
		fmt.Println(fmt.Errorf("Id is empty, get all users"))
	} else if len(keys[0]) > 0 {
		id = keys[0]
	}

	friendList, err := Store.GetFriendship(id)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Convert the `personList` variable to JSON
	friendListBytes, err := json.Marshal(friendList)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Write JSON list of persons to response
	w.Write(friendListBytes)
}

//CreateFriendshipHandler ...
func CreateFriendshipHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
	defer r.Body.Close()

	friendship := models.Friendship{}
	err = json.Unmarshal(body, &friendship)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	id, err := Store.CreateFriendship(&friendship)
	friendship.ID = id

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}

//DeleteFriendshipHandler ...
func DeleteFriendshipHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println(fmt.Errorf("Value does not exist"))
		return
	}

	id := keys[0]

	err := Store.DelFriendship(id)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}
