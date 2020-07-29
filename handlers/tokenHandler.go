package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func getByKey(r *http.Request, key string) string {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		fmt.Println(fmt.Errorf("Value does not exist"))
		return ""
	}

	return keys[0]
}

//CreateTokenHandler ...
func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	username := getByKey(r, "username")

	if username == "" {
		fmt.Println(fmt.Errorf("Empty username"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty username"))

		return
	}

	token, err := Redis.CreateToken(username, 60*24*time.Hour)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(token))
}

//CheckTokenHandler ...
func CheckTokenHandler(w http.ResponseWriter, r *http.Request) {
	username := getByKey(r, "username")
	token := getByKey(r, "token")

	if username == "" || token == "" {
		fmt.Println(fmt.Errorf("Empty username or token"))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Empty username or token"))
		return
	}

	ok, err := Redis.CheckToken(username, token)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad token"))
		return
	}
	w.Write([]byte("OK"))
}
