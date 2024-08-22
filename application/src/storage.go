package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

var dataStore sync.Map

func AddDataToStore(user UserTest) {
	dataStore.Store(user.UserID, user)
}

func GetDataFromStore(userId string) UserTest {
	if val, ok := dataStore.Load(userId); ok {
		user, ok := val.(UserTest)
		if !ok {
			panic("we only store UserTest in the map. nothing else should exist in it")
		}
		return user
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
