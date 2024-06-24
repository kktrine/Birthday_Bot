package handlers

import (
	"birthday_bot/internal/model"
	"birthday_bot/internal/storage"
	"encoding/json"
	"net/http"
)

func SubscribeHandler(db *storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data model.Subscribe
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil || data.SubscribeTo == nil || data.Id == nil || len(*data.SubscribeTo) == 0 {
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, "Some recovered fields are empty", http.StatusBadRequest)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = db.Subscribe(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
