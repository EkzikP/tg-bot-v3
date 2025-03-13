package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	TelegramBotToken string            `json:"telegram_bot_token"`
	ApiKey           string            `json:"api_key"`
	Host             string            `json:"host"`
	PhoneEngineer    map[string]string `json:"phone_engineer"`
}

func Read() Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		log.Panic("Error decoding config:", err)
	}
	return cfg
}
