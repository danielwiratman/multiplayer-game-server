package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"server/internal/env"
	"server/internal/models"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	rdb *redis.Client
	db  *pgx.Conn
}

func NewService(db *pgx.Conn, rdb *redis.Client) *Service {
	return &Service{
		rdb: rdb,
		db:  db,
	}
}

func (s *Service) Register(req models.RegisterRequest) error {
	req.Email = strings.ToLower(req.Email)

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return fmt.Errorf("bcrypt generate: %w", err)
	}
	slog.Debug("hashing password with bcrypt", "password", req.Password, "hash", string(passHash))

	if _, err := s.db.Exec(context.Background(),
		"INSERT INTO accounts (email, password_hash) VALUES ($1, $2)",
		req.Email, string(passHash)); err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("user already exists")
		}

		return fmt.Errorf("postgres insert: %w", err)
	}

	return nil
}

func (s *Service) Login(req models.LoginRequest) (*models.LoginResult, error) {
	type userType struct {
		ID           int
		PasswordHash string
	}

	user := userType{}
	if err := pgxscan.Get(context.Background(), s.db, &user,
		"SELECT id, password_hash FROM accounts WHERE email = $1",
		req.Email,
	); err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.NoDataFound {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("postgres select: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("wrong password")
	}

	newAccessToken, err := genAccessToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("jwt generate: %w", err)
	}

	newRefreshToken, refreshTokenJwtID, err := genRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("refresh token generate: %w", err)
	}

	cacheKey := "refresh:" + refreshTokenJwtID
	cacheVal := models.RefreshTokenCacheVal{
		UserID:    user.ID,
		ExpiresIn: time.Now().Add(30 * 24 * time.Hour),
	}
	cacheValBytes, err := json.Marshal(cacheVal)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	if err := s.rdb.Set(context.Background(), cacheKey, cacheValBytes, 30*24*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("redis set: %w", err)
	}

	return &models.LoginResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		UserID:       user.ID,
		ExpiresIn:    time.Now().Add(time.Duration(env.C.AccessTokenExpirationSeconds) * time.Second),
	}, nil
}

func genAccessToken(uid int) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "multiplayer-game-server",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(env.C.AccessTokenExpirationSeconds) * time.Second)),
		Subject:   fmt.Sprintf("%d", uid),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(env.C.JWTSecret))
}

func genRefreshToken(uid int) (string, string, error) {
	jti := uuid.NewString()

	claims := jwt.RegisteredClaims{
		Issuer:    "multiplayer-game-server",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(env.C.RefreshTokenExpirationSeconds) * time.Second)),
		Subject:   fmt.Sprintf("%d", uid),
		ID:        jti,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(env.C.JWTSecret))
	if err != nil {
		return "", "", fmt.Errorf("jwt generate: %w", err)
	}

	return token, jti, nil
}
