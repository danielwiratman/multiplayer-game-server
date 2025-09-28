package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"server/internal/auth"
	"server/internal/env"
	"server/internal/player"
	"server/pkg/betools"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	l := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
	slog.SetDefault(l)

	slog.Info("env loading")
	if err := env.Load(); err != nil {
		panic("env load: " + err.Error())
	}
	slog.Info("env loaded")

	slog.Info("redis connecting")
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", env.C.ValkeyHost, env.C.ValkeyPort),
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic("redis ping: " + err.Error())
	}
	slog.Info("redis connected")

	slog.Info("postgres connecting")
	db, err := pgx.Connect(context.Background(), fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.C.DBHost, env.C.DBPort, env.C.DBUser, env.C.DBPass, env.C.DBName))
	if err != nil {
		panic("postgres connect: " + err.Error())
	}

	if err := db.Ping(context.Background()); err != nil {
		panic("postgres ping: " + err.Error())
	}
	slog.Info("postgres connected")

	r := chi.NewRouter()

	authService := auth.NewService(db, rdb)
	authController := auth.NewController(authService)
	playerService := player.NewService(db, rdb)
	playerController := player.NewController(playerService)

	router := betools.NewRouter(
		authController,
		playerController,
	)

	r.Route("/", router.Route)

	slog.Info("server started", "port", env.C.ListenPort)
	http.ListenAndServe(":"+env.C.ListenPort, r)
}
