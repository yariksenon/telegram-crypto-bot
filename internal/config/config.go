package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	BotToken string
	ChatID   int64
}

func LoadConfig(logger *slog.Logger) (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")

	if botToken == "" || chatIDStr == "" {
		logger.Error("Необходимые переменные окружения отсутствуют", "BOT_TOKEN", botToken, "CHAT_ID", chatIDStr)
		return nil, errors.New("необходимые переменные окружения отсутствуют")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		logger.Error("Ошибка преобразования CHAT_ID в int64", "error", err)
		return nil, err
	}

	return &Config{
		BotToken: botToken,
		ChatID:   chatID,
	}, nil
}
