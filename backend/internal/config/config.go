package config

import "os"

type Config struct {
    MongoURI      string
    DBName        string
    JWTSecret     string
    ServerPort    string
    GoogleClientID     string
    GoogleClientSecret string
    GoogleRedirectURL  string
	OpenAIKey string
}

func New() *Config {
    return &Config{
        MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
        DBName:        getEnv("DB_NAME", "auth_db"),
        JWTSecret:     getEnv("JWT_SECRET", "secret"),
        ServerPort:    getEnv("SERVER_PORT", "8080"),
        GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		OpenAIKey: os.Getenv("OPENAI_API_KEY"),
    }
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}
