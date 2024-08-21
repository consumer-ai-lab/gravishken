package main

import (
	"encoding/json"
	"net/http"
)

type DataStore map[string]UserTest

var dataStore DataStore = make(DataStore)

func AddDataToStore(user UserTest) {
	dataStore[user.UserID] = user
}

func GetDataFromStore(userId string) UserTest {
	if val, ok := dataStore[userId]; ok {
		return val
	}
	return UserTest{}
}

func getDataFromStore(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	userData := GetDataFromStore(userId)
	if userData.UserID == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(userData)
	if err != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
		return
	}
}