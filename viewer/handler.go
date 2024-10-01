package viewer

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := DefaultViewer{}

		authorization := r.Header.Get("Authorization")
		if len(authorization) > 7 {
			token, _, err := new(jwt.Parser).ParseUnverified(authorization[7:], jwt.MapClaims{})
			if err != nil {
				log.Print("parse unverified jwt: ", err)
			} else {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					user.sub = claims["sub"].(string)
					user.jti = claims["jti"].(string)
					user.scopes = strings.Split(claims["scope"].(string), " ")
				}
			}
		}

		r = r.WithContext(newContext(r.Context(), user))
		next.ServeHTTP(w, r)
	})
}
