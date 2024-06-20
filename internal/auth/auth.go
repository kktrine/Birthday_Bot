package auth

import (
	"birthday_bot/internal/model"
	"birthday_bot/internal/storage"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Auth struct {
	JWTKey []byte
	Db     *storage.Storage
}

type RegisterResponse struct {
	Id int `json:"id"`
}

func (h *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(user.Username) < 5 || len(user.Password) < 5 {
		http.Error(w, "минимальная длина логина и пароля - 5 символов", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	err = h.Db.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	//err = json.NewEncoder(w).Encode(map[string]int{"id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Auth) SignIn(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	password := user.Password
	user, err = h.Db.GetHashedPassword(user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	hashedPassword := user.Password
	chatId := user.ChatId
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
	}
	if user.ChatId != nil && chatId != nil && *user.ChatId != *chatId {
		http.Error(w, "Невозможно войти с чужого устройства", http.StatusUnauthorized)
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claim := &claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(h.JWTKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"token": tokenString, "id": user.Id})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
