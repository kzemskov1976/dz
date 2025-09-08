package middleware

import (
	"context"
	"demo/order/6-order-api-cart/configs"
	"demo/order/6-order-api-cart/pkg/jwt"
	"net/http"
	"strings"
)

type phone string

const (
	PhoneKey phone = "PhoneKey"
)

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func Auth(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeUnauthorized(w)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		isValid, data := jwt.NewJWT(config.Secret).Parse(token)
		if !isValid {
			writeUnauthorized(w)
			return
		}
		ctx := context.WithValue(r.Context(), PhoneKey, data.Phone)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
