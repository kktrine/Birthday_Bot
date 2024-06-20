package main

import (
	"birthday_bot/internal/auth"
	"birthday_bot/internal/handlers"
	"birthday_bot/internal/middleware"
	"birthday_bot/internal/storage"
	"birthday_bot/internal/tgBot"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
	"time"
)

var jwtKey = []byte("my_secret_key")

func RunServer(address string) {
	db := storage.New("host=localhost user=postgres password=postgres dbname=birthdays port=5432 sslmode=disable")
	authProcess := auth.Auth{JWTKey: jwtKey, Db: db}

	router := mux.NewRouter()
	router.HandleFunc("/sign_up", authProcess.Register).Methods("POST")
	router.HandleFunc("/sign_in", authProcess.SignIn).Methods("POST")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(jwtKey))
	protected.HandleFunc("/employees", handlers.GetEmployees(db)).Methods("GET")
	protected.HandleFunc("/info", handlers.PostInfo(db)).Methods("POST")

	log.Println("Server is running on " + address)
	log.Fatal(http.ListenAndServe(address, router))
}

func RunBot(address string) {
	log.Println("Starting tg bot")
	bot := tgBot.NewBot("http://" + address)
	bot.Start()
}

func main() {
	address := "localhost:8080"
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		RunServer(address)
	}()
	time.Sleep(3 * time.Second)
	go func() {
		defer wg.Done()
		RunBot(address)
	}()

	wg.Wait()
}
