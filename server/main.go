package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"server/internal/auth"
	"server/internal/env"
	"server/pkg/betools"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
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

	r := chi.NewRouter()

	authController := auth.NewController()

	router := betools.NewRouter(
		authController,
	)

	r.Route("/", router.Route)

	slog.Info("server started", "port", env.C.ListenPort)
	http.ListenAndServe(":"+env.C.ListenPort, r)
}
