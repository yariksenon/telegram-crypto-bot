package usecases

import (
	"context"
	"fmt"
	"strings"
	"tg-crypto-bot/internal/adapters"
	"tg-crypto-bot/internal/interfaces"
	"time"
)

type CryptoUsecase interface {
	SendCurrentPrice(ctx context.Context, symbol string) error
	HandleUserMessage(ctx context.Context, message string, chatID int64) error
}

type cryptoUsecase struct {
	priceProvider interfaces.PriceProvider
	telegram      *adapters.TelegramAdapter
}

func NewCryptoUsecase(pp interfaces.PriceProvider, tg *adapters.TelegramAdapter) CryptoUsecase {
	return &cryptoUsecase{
		priceProvider: pp,
		telegram:      tg,
	}
}

func (c *cryptoUsecase) SendCurrentPrice(ctx context.Context, symbol string) error {
	price, err := c.priceProvider.GetPrice(symbol)
	if err != nil {
		return fmt.Errorf("ошибка получения цены от Binance для %s: %w", symbol, err)
	}

	message := fmt.Sprintf("Текущая цена %s: %.4f", symbol, price)

	if err := c.telegram.SendMessage(ctx, message); err != nil {
		return fmt.Errorf("ошибка отправки сообщения в Telegram: %w", err)
	}
	return nil
}

func (c *cryptoUsecase) HandleUserMessage(ctx context.Context, message string, chatID int64) error {
	// Убираем лишние пробелы и приводим к верхнему регистру
	symbol := strings.TrimSpace(strings.ToUpper(message))

	// Проверяем, что сообщение похоже на торговую	 пару
	if !strings.Contains(symbol, "USDT") && !strings.Contains(symbol, "BTC") && !strings.Contains(symbol, "ETH") {
		helpMsg := `Отправьте мне символ криптовалюты для получения текущей цены.

Примеры:
• BTCUSDT - цена Bitcoin
• ETHUSDT - цена Ethereum  
• ADAUSDT - цена Cardano
• DOTUSDT - цена Polkadot
`

		return c.telegram.SendMessageToChat(ctx, helpMsg, chatID)
	}

	// Если пользователь отправил только символ без USDT, добавляем USDT
	if !strings.HasSuffix(symbol, "USDT") && len(symbol) <= 5 {
		symbol = symbol + "USDT"
	}

	// Получаем цену
	price, err := c.priceProvider.GetPrice(symbol)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось получить цену для %s: %v", symbol, err)
		return c.telegram.SendMessageToChat(ctx, errorMsg, chatID)
	}

	// Форматируем сообщение с ценой
	priceMsg := fmt.Sprintf(
		"💰 %s\n💵 Цена: %.2f USDT\n⌚ Время: %s\n🏦 Источник: Binance",
		symbol,
		price,
		time.Now().Format("15:04:05"),
	)

	return c.telegram.SendMessageToChat(ctx, priceMsg, chatID)
}
