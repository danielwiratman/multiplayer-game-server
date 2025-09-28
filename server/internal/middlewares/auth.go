package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"server/internal/env"
	"server/pkg/betools"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken, ok := getBearer(r.Header.Get("Authorization"))
		if !ok {
			slog.Error("auth middleware", "error", "missing authorization header")
			betools.SendErrorResponse(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		token, err := jwt.ParseWithClaims(bearerToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(env.C.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			slog.Error("auth middleware", "error", err.Error())
			betools.SendErrorResponse(w, http.StatusUnauthorized, "token invalid")
			return
		}

		uid, err := strconv.Atoi(token.Claims.(*jwt.RegisteredClaims).Subject)
		if err != nil {
			slog.Error("auth middleware", "error", err.Error())
			betools.SendErrorResponse(w, http.StatusUnauthorized, "token invalid")
			return
		}

		next.ServeHTTP(w, betools.SetContext(r, betools.CtxKeyAuth, uid))
	})
}

func getBearer(header string) (string, bool) {
	if !strings.HasPrefix(header, "Bearer ") {
		return "", false
	}
	return strings.TrimPrefix(header, "Bearer "), true
}
