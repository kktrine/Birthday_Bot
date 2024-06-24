package handlers

import (
	"birthday_bot/internal/model"
	"birthday_bot/internal/storage"
	"encoding/json"
	"net/http"
	"time"
)

func PostInfoHandler(db *storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var person model.Employee
		err := json.NewDecoder(r.Body).Decode(&person)
		if err != nil || person.Birth == nil || person.Birth.After(time.Now()) {
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, "wrong field birth", http.StatusBadRequest)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = db.AddInfo(person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
