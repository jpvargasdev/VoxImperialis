package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// AppConfig holds all runtime configuration for Vox Imperialis.
type AppConfig struct {
	JID             string
	Password        string
	Server          string
	AllowedUsers    []string
	AllowedServices []string
	TLSSkipVerify   bool
}

var appConfig AppConfig

// Load reads configuration from .env file then environment variables.
// Required: XMPP_JID, XMPP_PASSWORD, XMPP_SERVER, ALLOWED_USERS
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	appConfig = AppConfig{
		JID:             mustEnv("XMPP_JID"),
		Password:        mustEnv("XMPP_PASSWORD"),
		Server:          mustEnv("XMPP_SERVER"),
		AllowedUsers:    splitList(mustEnv("ALLOWED_USERS")),
		AllowedServices: splitList(getEnv("ALLOWED_SERVICES", "nginx,caddy,tailscaled")),
		TLSSkipVerify:   getEnv("XMPP_TLS_SKIP_VERIFY", "false") == "true",
	}

	if len(appConfig.AllowedUsers) == 0 {
		log.Fatal("ALLOWED_USERS must contain at least one JID")
	}

	log.Printf("config loaded: server=%s jid=%s users=%d services=%d",
		appConfig.Server, appConfig.JID,
		len(appConfig.AllowedUsers), len(appConfig.AllowedServices),
	)
}

// Get returns the loaded configuration.
func Get() AppConfig {
	return appConfig
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env variable %q is not set", key)
	}
	return v
}

func getEnv(key, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return defaultVal
}

// splitList splits a comma-separated string into a trimmed, non-empty slice.
func splitList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
