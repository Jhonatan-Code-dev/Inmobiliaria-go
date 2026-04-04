package env

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppEnv       string `envconfig:"APP_ENV" default:"development"`
	Port         string `envconfig:"PORT" default:"4000"`
	DSN          string `envconfig:"BASE_DATOS_1" required:"true"`
	JWTSecret    string `envconfig:"JWT_SECRET" required:"true"`
}

var (
	cfg  *Config
	once sync.Once
)

func NewConfig() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		cfg = &Config{}
		if err := envconfig.Process("", cfg); err != nil {
			log.Fatalf("error al cargar variables de entorno: %v", err)
		}
	})

	return cfg
}
