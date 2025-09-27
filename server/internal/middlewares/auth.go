package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"server/internal/env"
	"server/internal/models"
	"server/pkg/betools"
	"strconv"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type Middlewares struct {
	rdb *redis.Client
	db  *pgx.Conn
}

func NewMiddlewares(rdb *redis.Client, db *pgx.Conn) *Middlewares {
	return &Middlewares{
		rdb: rdb,
		db:  db,
	}
}

func (m *Middlewares) AuthMiddleware() betools.Middleware {
	return func(next http.Handler) http.Handler {
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

			account := models.Account{}

			accountInfoJson, err := m.rdb.Get(context.Background(), fmt.Sprintf("account:%d", uid)).Result()
			if err == redis.Nil {
				if err := pgxscan.Get(context.Background(), m.db, &account,
					"SELECT first_name, last_name, email FROM accounts WHERE id = $1",
					uid); err != nil {
					slog.Error("auth middleware", "error", err.Error())
					betools.SendErrorResponse(w, http.StatusInternalServerError, "token invalid")
					return
				}
				accountInfoJsonFromDB, err := json.Marshal(account)
				if err != nil {
					slog.Error("auth middleware", "error", err.Error())
					betools.SendErrorResponse(w, http.StatusInternalServerError, "token invalid")
					return
				}
				if err := m.rdb.Set(context.Background(), fmt.Sprintf("account:%d", uid), string(accountInfoJsonFromDB), time.Duration(env.C.AccountInfoCacheSeconds)*time.Second).Err(); err != nil {
					slog.Error("auth middleware", "error", err.Error())
					betools.SendErrorResponse(w, http.StatusInternalServerError, "token invalid")
					return
				}
			} else if err != nil {
				slog.Error("auth middleware", "error", err.Error())
				betools.SendErrorResponse(w, http.StatusInternalServerError, "token invalid")
				return
			} else {
				if err := json.Unmarshal([]byte(accountInfoJson), &account); err != nil {
					slog.Error("auth middleware", "error", err.Error())
					betools.SendErrorResponse(w, http.StatusInternalServerError, "token invalid")
					return
				}
			}

			slog.Debug("auth middleware", "account", account)

			next.ServeHTTP(w, betools.SetContext(r, betools.CtxKeyAuth, account))
		})
	}
}

func getBearer(header string) (string, bool) {
	if !strings.HasPrefix(header, "Bearer ") {
		return "", false
	}
	return strings.TrimPrefix(header, "Bearer "), true
}
