package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	ApiKey map[string]string
}

func LoadConfig() *config {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API key not found in environment variables")
	}
	return &config{
		ApiKey: make(map[string]string),
	}
}
func (c *config) GetApiKey(service string) string {
	if key, exists := c.ApiKey[service]; exists {
		return key
	}
	key := os.Getenv(service)
	if key == "" {
		log.Fatalf("API key for service %s not found in environment variables", service)
	}
	c.ApiKey[service] = key
	return key
}

func (c *config) SetApiKey(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("API key %s not found in environment variables", key)
	}
	return value
}
