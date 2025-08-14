package interfaces

import "context"

type PriceProvider interface {
	GetPrice(symbol string) (float64, error)
}

type TelegramSender interface {
	SendMessage(ctx context.Context, message string) error
}

type MessageHandler interface {
	HandleMessage(ctx context.Context, message string, chatID int64) error
}

type BotRunner interface {
	Start(ctx context.Context) error
}
