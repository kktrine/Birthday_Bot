package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func AuthMiddleware(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			claim := &claims{}
			tkn, err := jwt.ParseWithClaims(authHeader, claim, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil {
				if errors.Is(err, jwt.ErrSignatureInvalid) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if !tkn.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
