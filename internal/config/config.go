package config

import (
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
	Env       string
	LogLevel  string
	Timezone  string
	JWTSecret string
}

type PostgresConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (c *PostgresConfig) DSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.Database + "?sslmode=" + c.SSLMode
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type WorkerConfig struct {
	HoldReleaseInterval   time.Duration
	BookingExpireInterval time.Duration
	HoldDurationMinutes   int
	BookingExpireMinutes  int
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/gobooking")

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

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
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
