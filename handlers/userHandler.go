package handlers

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//GetUsersHandler ...
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve people from postgresql database using our `store` interface variable's
	// `func (*dbstore) GetUser` pointer receiver method defined in `store.go` file
	id := ""
	keys, ok := r.URL.Query()["id"]
	if !ok {
		fmt.Println(fmt.Errorf("Id is empty, get all users"))
	} else if len(keys[0]) > 0 {
		id = keys[0]
	}

	userList, err := Store.GetUsers(id)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Convert the `personList` variable to JSON
	userListBytes, err := json.Marshal(userList)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Write JSON list of persons to response
	w.Write(userListBytes)
}

//CreateUserHandler ...
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
	defer r.Body.Close()

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	id, err := Store.CreateUser(&user)
	user.ID = id
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}

//UpdateUserHandler ...
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	err = Store.UpdUser(&user, id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}

//DeleteUserHandler ...
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println(fmt.Errorf("Value does not exist"))
		return
	}

	id := keys[0]

	err := Store.DelUser(id)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
}
