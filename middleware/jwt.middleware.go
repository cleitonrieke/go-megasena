package middleware

import (
	"backlotofacil/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/context"
	"gorm.io/gorm"
)

type HandlerMiddleware struct {
	DB *gorm.DB
}

func (h HandlerMiddleware) ValidarJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read basic auth information
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}
		id := claims.(jwt.MapClaims)["sub"].(float64)

		context.Set(r, "userId", id)

		next.ServeHTTP(w, r)
	}
}
