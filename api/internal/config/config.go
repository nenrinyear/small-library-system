package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           int
	DBHost         string
	DBPort         int
	DBName         string
	DBUser         string
	DBPassword     string
	DBTLS          string
	DBParseTime    bool
	AdminToken     string
	AllowedOrigins []string
	HistoryLimit   int
}

func Load() Config {
	cfg := Config{
		Port:           intEnv("API_PORT", 8080),
		DBHost:         stringEnv("MYSQL_HOST", "localhost"),
		DBPort:         intEnv("MYSQL_PORT", 3306),
		DBName:         stringEnv("MYSQL_DATABASE", "small_library_system"),
		DBUser:         stringEnv("MYSQL_USER", "sls_user"),
		DBPassword:     stringEnv("MYSQL_PASSWORD", ""),
		DBTLS:          stringEnv("MYSQL_TLS", "false"),
		DBParseTime:    boolEnv("MYSQL_PARSE_TIME", true),
		AdminToken:     stringEnv("ADMIN_API_TOKEN", ""),
		AllowedOrigins: csvEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
		HistoryLimit:   intEnv("MAX_HISTORY_COUNT", 10),
	}

	if cfg.HistoryLimit <= 0 {
		cfg.HistoryLimit = 10
	}

	return cfg
}

func (c Config) MySQLDSN() string {
	return c.mysqlDSN(false)
}

func (c Config) MySQLMigrationDSN() string {
	return c.mysqlDSN(true)
}

func (c Config) mysqlDSN(multiStatements bool) string {
	parseTime := "false"
	if c.DBParseTime {
		parseTime = "true"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=%s&loc=Asia%%2FTokyo&tls=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		parseTime,
		c.DBTLS,
	)
	if multiStatements {
		dsn += "&multiStatements=true"
	}
	return dsn
}

func stringEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func intEnv(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func boolEnv(key string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func csvEnv(key, fallback string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		raw = fallback
	}

	parts := strings.Split(raw, ",")
	items := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			items = append(items, v)
		}
	}

	if len(items) == 0 {
		return []string{"http://localhost:3000"}
	}

	return items
}
