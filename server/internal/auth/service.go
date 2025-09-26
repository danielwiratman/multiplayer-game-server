package auth

import (
	"context"
	"fmt"
	"log/slog"
	"server/internal/models"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golang-jwt/jwt/v5"
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

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MaxCost)
	if err != nil {
		return fmt.Errorf("bcrypt generate: %w", err)
	}
	slog.Debug("hashing password with bcrypt", "password", req.Password, "hash", string(passHash))

	if _, err := s.db.Exec(context.Background(),
		"INSERT INTO users (email, password_hash) VALUES ($1, $2)",
		req.Email, string(passHash)); err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("user already exists")
		}

		return fmt.Errorf("postgres insert: %w", err)
	}

	return nil
}

func (s *Service) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	type userType struct {
		ID           int
		PasswordHash string
	}

	user := userType{}
	if err := pgxscan.Select(context.Background(), s.db, &user,
		"SELECT id, password_hash FROM accounts WHERE email = $1",
		req.Email,
	); err != nil {
		return nil, fmt.Errorf("postgres select: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("bcrypt compare: %w", err)
	}

	newAccessToken, err := s.genJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("jwt generate: %w", err)
	}

	return &models.LoginResponse{
		AccessToken: newAccessToken,
	}, nil
}

func (s *Service) genJWT(uid int) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "multiplayer-game-server",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		Subject:   fmt.Sprintf("%d", uid),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
}
