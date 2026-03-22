package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	JWT      JWTConfig
	Storage  StorageConfig
	SMTP     SMTPConfig
}

type AppConfig struct {
	Name string
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
	RefreshHours    int
}

type StorageConfig struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	Region    string
}

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	FromName  string
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "koperasi-api"),
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "koperasi"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			GroupID: getEnv("KAFKA_GROUP_ID", "koperasi-api"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "change-me-in-production"),
			ExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
			RefreshHours:    getEnvInt("JWT_REFRESH_HOURS", 168),
		},
		Storage: StorageConfig{
			Endpoint:  getEnv("STORAGE_ENDPOINT", ""),
			Bucket:    getEnv("STORAGE_BUCKET", "koperasi"),
			AccessKey: getEnv("STORAGE_ACCESS_KEY", ""),
			SecretKey: getEnv("STORAGE_SECRET_KEY", ""),
			Region:    getEnv("STORAGE_REGION", "us-east-1"),
		},
		SMTP: SMTPConfig{
			Host:      getEnv("SMTP_HOST", "smtp.mailtrap.io"),
			Port:      getEnvInt("SMTP_PORT", 2525),
			Username:  getEnv("SMTP_USERNAME", ""),
			Password:  getEnv("SMTP_PASSWORD", ""),
			FromEmail: getEnv("SMTP_FROM_EMAIL", "noreply@koperasi.local"),
			FromName:  getEnv("SMTP_FROM_NAME", "Koperasi System"),
		},
	}
}

func (c *DatabaseConfig) DSN() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Warning: invalid int for %s: %s, using fallback %d", key, value, fallback)
			return fallback
		}
		return i
	}
	return fallback
}
