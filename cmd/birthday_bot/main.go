package main

import (
	"birthday_bot/internal/auth"
	"birthday_bot/internal/handlers"
	"birthday_bot/internal/middleware"
	"birthday_bot/internal/storage"
	"birthday_bot/internal/tgBot"
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var jwtKey = []byte("my_secret_key")

func RunServer(address string, db *storage.Storage) {

	authProcess := auth.Auth{JWTKey: jwtKey, Db: db}

	router := mux.NewRouter()
	router.HandleFunc("/sign_up", authProcess.Register).Methods("POST")
	router.HandleFunc("/sign_in", authProcess.SignIn).Methods("POST")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(jwtKey))
	protected.HandleFunc("/employees", handlers.GetEmployeesHandler(db)).Methods("GET")
	protected.HandleFunc("/info", handlers.PostInfoHandler(db)).Methods("POST")
	protected.HandleFunc("/subscribe", handlers.SubscribeHandler(db)).Methods("POST")
	protected.HandleFunc("/unsubscribe", handlers.UnsubscribeHandler(db)).Methods("POST")

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	// Канал для уведомления о завершении работы сервера
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Server is running on " + address)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-stop

	log.Println("Server is shutting down...")
	_ = db.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

func RunBot(address string, db *storage.Storage) {
	log.Println("Starting tg bot")
	bot := tgBot.NewBot("http://" + address)
	bot.Start(db)
}

func main() {
	err := godotenv.Load("./env/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	address := "localhost:" + os.Getenv("PORT")
	db := storage.New(os.Getenv("POSTGRES"))
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		RunServer(address, db)
	}()
	time.Sleep(3 * time.Second)
	go func() {
		RunBot(address, db)
	}()

	wg.Wait()
}
