package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App     AppConfig  `yaml:"app"`
		HTTP    HTTPConfig `yaml:"http"`
		Log     LogConfig  `yaml:"log"`
		DB      DBConfig   `yaml:"db"`
		Redis   RedisConfig `yaml:"redis"`
		JWT     JWTConfig  `yaml:"jwt"`
		GodUser GodUserConfig `yaml:"god_user"`
	}

	AppConfig struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		Env     string `env-required:"true" yaml:"env" env:"APP_ENV"` // local, dev, prod
	}

	HTTPConfig struct {
		Port            string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		ReadTimeout     time.Duration `env-required:"true" yaml:"read_timeout" env:"HTTP_READ_TIMEOUT"`
		WriteTimeout    time.Duration `env-required:"true" yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT"`
		ShutdownTimeout time.Duration `env-required:"true" yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT"`
	}

	LogConfig struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"` // debug, info, warn, error
	}

	DBConfig struct {
		URL             string        `env-required:"true" yaml:"url" env:"DATABASE_URL"`
		MaxConns        int32         `env-default:"10" yaml:"max_conns" env:"DB_MAX_CONNS"`
		MinConns        int32         `env-default:"2" yaml:"min_conns" env:"DB_MIN_CONNS"`
		MaxConnIdleTime time.Duration `env-default:"15m" yaml:"max_conn_idle_time" env:"DB_MAX_CONN_IDLE_TIME"`
		MaxConnLifeTime time.Duration `env-default:"1h" yaml:"max_conn_lifetime" env:"DB_MAX_CONN_LIFETIME"`
	}

	RedisConfig struct {
		URL string `env-required:"true" yaml:"url" env:"REDIS_URL"`
	}

	JWTConfig struct {
		Secret            string        `env-required:"true" yaml:"secret" env:"JWT_SECRET"`
		AccessExpiration  time.Duration `env-required:"true" yaml:"access_expiration" env:"JWT_ACCESS_EXPIRATION"`
		RefreshExpiration time.Duration `env-required:"true" yaml:"refresh_expiration" env:"JWT_REFRESH_EXPIRATION"`
	}

	GodUserConfig struct {
		Email    string `env-default:"god@admin.com" yaml:"email" env:"GOD_USER_EMAIL"`
		Password string `env-default:"god-admin-password" yaml:"password" env:"GOD_USER_PASSWORD"`
		Name     string `env-default:"God" yaml:"name" env:"GOD_USER_NAME"`
	}
)

// LoadConfig loads config from environment variables and/or optional .env file
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// If .env file exists, cleanenv will read it. Otherwise it falls back to system env vars.
	var err error
	if _, errStat := os.Stat(".env"); errStat == nil {
		err = cleanenv.ReadConfig(".env", cfg)
	} else {
		err = cleanenv.ReadEnv(cfg)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}
