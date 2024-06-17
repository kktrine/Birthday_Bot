package main

import (
	"birthday_bot/internal/auth"
	"birthday_bot/internal/handlers"
	"birthday_bot/internal/middleware"
	"birthday_bot/internal/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var jwtKey = []byte("my_secret_key")

func main() {

	db := storage.New("host=localhost user=postgres password=postgres dbname=birthdays port=5432 sslmode=disable")
	authProcess := auth.Auth{JWTKey: jwtKey, Db: db}

	router := mux.NewRouter()
	router.HandleFunc("/sign_up", authProcess.Register).Methods("POST")
	router.HandleFunc("/sign_in", authProcess.SignIn).Methods("POST")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(jwtKey))
	protected.HandleFunc("/employees", handlers.GetEmployees(db)).Methods("GET")

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
