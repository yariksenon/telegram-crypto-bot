package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tg-crypto-bot/internal/adapters"
	"tg-crypto-bot/internal/config"
	"tg-crypto-bot/internal/usecases"

	"github.com/joho/godotenv"
)

func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return logger
}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("Файл .env не найден или не загружен")
	}

	//Настройка логгера
	logger := setupLogger()

	//Загрузка конфигурации
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Error("Ошибка загрузки конфигурации", "error", err)
		os.Exit(1)
	}

	//Инициализация адаптеров
	binanceProvider := adapters.NewBinancePriceProvider()
	telegramAdapter, err := adapters.NewTelegramAdapter(cfg.BotToken, cfg.ChatID)
	if err != nil {
		logger.Error("Не удалось инициализировать TelegramAdapter", "error", err)
		os.Exit(1)
	}

	//Инициализация сервисов
	uc := usecases.NewCryptoUsecase(binanceProvider, telegramAdapter)

	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработчик сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем бота в горутине
	go func() {
		logger.Info("Бот запущен и ожидает сообщения...")
		if err := telegramAdapter.Start(ctx, func(message string, chatID int64) error {
			return uc.HandleUserMessage(ctx, message, chatID)
		}); err != nil {
			logger.Error("Ошибка работы бота", "error", err)
		}
	}()

	// Ожидаем сигнал завершения
	<-sigChan
	logger.Info("Получен сигнал завершения, останавливаем бота...")
	cancel()
	logger.Info("Бот остановлен")
}
