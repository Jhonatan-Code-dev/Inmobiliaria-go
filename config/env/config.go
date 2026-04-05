package env

import (
	"log"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppEnv       string `envconfig:"APP_ENV" default:"development"`
	Port         string `envconfig:"PORT" default:"4000"`
	DSN          string `envconfig:"BASE_DATOS_1" required:"true"`
	JWTSecret    string `envconfig:"JWT_SECRET" required:"true"`
	AllowedOrigins string `envconfig:"ALLOWED_ORIGINS" default:"*"`
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

		cfg.Port = normalizarPuerto(cfg.Port)
	})

	return cfg
}

func normalizarPuerto(port string) string {
	port = strings.TrimSpace(port)
	port = strings.TrimPrefix(port, ":")

	if port == "" {
		return "4000"
	}

	return port
}
