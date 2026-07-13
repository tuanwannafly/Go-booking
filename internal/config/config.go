package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Server   ServerConfig
	Worker   WorkerConfig
}

type AppConfig struct {
	Env       string `mapstructure:"env"`
	LogLevel  string `mapstructure:"log_level"`
	Timezone  string `mapstructure:"timezone"`
	JWTSecret string `mapstructure:"jwt_secret"`
}

type PostgresConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

func (c *PostgresConfig) DSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.Database + "?sslmode=" + c.SSLMode
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type WorkerConfig struct {
	HoldReleaseInterval   time.Duration `mapstructure:"hold_release_interval"`
	BookingExpireInterval time.Duration `mapstructure:"booking_expire_interval"`
	HoldDurationMinutes   int           `mapstructure:"hold_duration_minutes"`
	BookingExpireMinutes  int           `mapstructure:"booking_expire_minutes"`
}

func Load() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	bindEnv("app.env", "APP_ENV")
	bindEnv("app.log_level", "LOG_LEVEL")
	bindEnv("app.timezone", "APP_TIMEZONE")
	bindEnv("app.jwt_secret", "JWT_SECRET")
	bindEnv("postgres.host", "DB_HOST")
	bindEnv("postgres.port", "DB_PORT")
	bindEnv("postgres.user", "DB_USER")
	bindEnv("postgres.password", "DB_PASSWORD")
	bindEnv("postgres.database", "DB_NAME")
	bindEnv("postgres.sslmode", "DB_SSLMODE")
	bindEnv("redis.host", "REDIS_HOST")
	bindEnv("redis.port", "REDIS_PORT")
	bindEnv("redis.password", "REDIS_PASSWORD")
	bindEnv("redis.db", "REDIS_DB")
	bindEnv("server.host", "SERVER_HOST")
	bindEnv("server.port", "SERVER_PORT")
	bindEnv("server.read_timeout", "SERVER_READ_TIMEOUT")
	bindEnv("server.write_timeout", "SERVER_WRITE_TIMEOUT")
	bindEnv("server.idle_timeout", "SERVER_IDLE_TIMEOUT")
	bindEnv("worker.hold_release_interval", "WORKER_HOLD_RELEASE_INTERVAL")
	bindEnv("worker.booking_expire_interval", "WORKER_BOOKING_EXPIRE_INTERVAL")
	bindEnv("worker.hold_duration_minutes", "WORKER_HOLD_DURATION_MINUTES")
	bindEnv("worker.booking_expire_minutes", "WORKER_BOOKING_EXPIRE_MINUTES")

	setDefaults()

	if err := loadDotEnv(".env"); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// loadDotEnv reads a .env file and exports each key into the process
// environment so viper's BindEnv bindings (DB_USER -> postgres.user, etc.)
// resolve correctly. Process env vars already set take precedence.
func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"'`)

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
	return scanner.Err()
}

func bindEnv(key, env string) {
	_ = viper.BindEnv(key, env)
}

func setDefaults() {
	// App
	viper.SetDefault("app.env", "development")
	viper.SetDefault("app.log_level", "debug")
	viper.SetDefault("app.timezone", "Asia/Ho_Chi_Minh")
	viper.SetDefault("app.jwt_secret", "dev-secret-change-in-production")

	// Postgres
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.user", "postgres")
	viper.SetDefault("postgres.password", "postgres")
	viper.SetDefault("postgres.database", "gobooking")
	viper.SetDefault("postgres.sslmode", "disable")
	viper.SetDefault("postgres.max_open_conns", 25)
	viper.SetDefault("postgres.max_idle_conns", 5)
	viper.SetDefault("postgres.conn_max_lifetime", "5m")

	// Redis
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Server
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.idle_timeout", "60s")

	// Worker
	viper.SetDefault("worker.hold_release_interval", "30s")
	viper.SetDefault("worker.booking_expire_interval", "60s")
	viper.SetDefault("worker.hold_duration_minutes", 10)
	viper.SetDefault("worker.booking_expire_minutes", 10)
}
