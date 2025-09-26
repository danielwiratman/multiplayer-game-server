package env

import "github.com/caarlos0/env/v10"

type Env struct {
	ValkeyHost string `env:"VALKEY_HOST"`
	ValkeyPort string `env:"VALKEY_PORT"`

	ListenPort string `env:"LISTEN_PORT" envDefault:"8080"`

	DBHost string `env:"DB_HOST"`
	DBPort string `env:"DB_PORT"`
	DBUser string `env:"DB_USER"`
	DBPass string `env:"DB_PASS"`
	DBName string `env:"DB_NAME"`
}

var C *Env

func Load() error {
	var c Env
	if err := env.ParseWithOptions(&c, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		return err
	}

	C = &c

	return nil
}
