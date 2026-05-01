package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port             string
	OAuthToken       string
	OAuthTokenSecret string
	ConsumerKey      string
	ConsumerSecret   string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	c := &Config{
		Port:             port,
		OAuthToken:       os.Getenv("OAUTH_TOKEN"),
		OAuthTokenSecret: os.Getenv("OAUTH_TOKEN_SECRET"),
		ConsumerKey:      os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:   os.Getenv("CONSUMER_SECRET"),
	}
	if c.OAuthToken == "" || c.OAuthTokenSecret == "" || c.ConsumerKey == "" || c.ConsumerSecret == "" {
		return nil, fmt.Errorf("missing Plurk OAuth env: OAUTH_TOKEN, OAUTH_TOKEN_SECRET, CONSUMER_KEY, CONSUMER_SECRET")
	}
	return c, nil
}
