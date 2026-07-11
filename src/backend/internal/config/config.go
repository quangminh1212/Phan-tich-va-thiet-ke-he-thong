package config

import "os"

type Config struct {
	DBDSN     string
	JWTSecret string
	Port      string
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() Config {
	return Config{
		DBDSN: getenv("DB_DSN",
			"host=localhost user=warehouse password=warehouse dbname=warehouse port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"),
		JWTSecret: getenv("JWT_SECRET", "dev-secret-change-me"),
		Port:      getenv("PORT", "8080"),
	}
}
