package env

import "github.com/caarlos0/env/v10"

type Env struct {
	IsProd bool `env:"IS_PROD" envDefault:"false"`

	ValkeyHost string `env:"VALKEY_HOST"`
	ValkeyPort string `env:"VALKEY_PORT"`

	ListenPort string `env:"LISTEN_PORT" envDefault:"8080"`

	DBHost string `env:"DB_HOST"`
	DBPort string `env:"DB_PORT"`
	DBUser string `env:"DB_USER"`
	DBPass string `env:"DB_PASS"`
	DBName string `env:"DB_NAME"`

	AccessTokenExpirationSeconds  int `env:"ACCESS_TOKEN_EXPIRATION_SECONDS" envDefault:"300"`
	RefreshTokenExpirationSeconds int `env:"REFRESH_TOKEN_EXPIRATION_SECONDS" envDefault:"604800"` // 7 days

	AccountInfoCacheSeconds int `env:"ACCOUNT_INFO_CACHE_SECONDS" envDefault:"60"`

	JWTSecret string `env:"JWT_SECRET"`
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
