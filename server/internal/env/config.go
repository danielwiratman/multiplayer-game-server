package env

import "github.com/caarlos0/env/v10"

type Env struct {
	ValkeyHost string `env:"VALKEY_HOST"`
	ValkeyPort string `env:"VALKEY_PORT"`

	ListenPort string `env:"LISTEN_PORT" envDefault:"8080"`
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
