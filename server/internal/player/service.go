package player

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/env"
	"server/internal/models"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
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

func (s *Service) GetInfo(uid int) (*models.Account, error) {
	account := models.Account{}

	accountInfoJson, err := s.rdb.Get(context.Background(), fmt.Sprintf("account:%d", uid)).Result()
	if err == redis.Nil {
		if err := pgxscan.Get(context.Background(), s.db, &account,
			"SELECT first_name, last_name, email FROM accounts WHERE id = $1",
			uid); err != nil {
			return nil, fmt.Errorf("postgres select: %w", err)
		}
		accountInfoJsonFromDB, err := json.Marshal(account)
		if err != nil {
			return nil, fmt.Errorf("json marshal: %w", err)
		}
		if err := s.rdb.Set(context.Background(), fmt.Sprintf("account:%d", uid), string(accountInfoJsonFromDB), time.Duration(env.C.AccountInfoCacheSeconds)*time.Second).Err(); err != nil {
			return nil, fmt.Errorf("redis set: %w", err)
		}
	} else if err != nil {
	} else {
		if err := json.Unmarshal([]byte(accountInfoJson), &account); err != nil {
			return nil, fmt.Errorf("json unmarshal: %w", err)
		}
	}

	return &account, nil
}

func (s *Service) UpdateInfo(uid int, req models.UpdatePlayerRequest) error {
	return nil
}
