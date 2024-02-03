package middleware

import (
	"context"
	"eau-de-go/pkg/jwt_util"
	"fmt"
	"net/http"
	"strings"
)

func getAccessTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header not found")
	}
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("Bearer token not found in authorization header")
	}

	accessTokenString := splitToken[1]
	return accessTokenString, nil
}

func JwtAuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		accessTokenString, err := getAccessTokenFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := jwt_util.DecodeToken(jwt_util.Access, accessTokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "jwt_claims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
