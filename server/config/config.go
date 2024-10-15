package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Host          string `json:"host"`
	Port          string `json:"port"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Database      string `json:"database"`
	RedisHost     string `json:"redis_host"`
	RedisPassword string `json:"redis_password"`
}

var conf *Config

func GetConfig() *Config {
	if conf != nil {
		return conf
	}

	// ... get configuration (from env vars or config file)
	conf = &Config{}

	configContent, err := os.ReadFile(".env")
	if err != nil {
		fmt.Println("Error reading .env file: ", err)
		os.Exit(1)
	}
	envVariables := strings.Split(string(configContent), "\n")
	envMap := make(map[string]string)
	for _, variable := range envVariables {
		if strings.Contains(variable, "=") {
			splitVariable := strings.Split(variable, "=")
			envMap[splitVariable[0]] = splitVariable[1]
		}
	}

	if err := json.Unmarshal([]byte(configContent), conf); err != nil {
		fmt.Println("Error unmarshalling config file: ", err)
		os.Exit(1)
	}

	// Apply environment overrides (if any)
	for k, v := range envMap {
		if strings.Contains(k, "host") {
			conf.Host = v
		} else if strings.Contains(k, "port") {
			conf.Port = v
		} else if strings.Contains(k, "user") {
			conf.User = v
		} else if strings.Contains(k, "password") {
			conf.Password = v
		} else if strings.Contains(k, "database") {
			conf.Database = v
		} else if strings.Contains(k, "redis_host") {
			conf.RedisHost = v
		} else if strings.Contains(k, "redis_password") {
			conf.RedisPassword = v
		}
	}

	return conf
}

// Initialize the config based on environment variables (e.g., .env file or similar)

func init() {
	conf = GetConfig()
}
