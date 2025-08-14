package adapters

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramAdapter struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramAdapter(botToken string, chatID int64) (*TelegramAdapter, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	return &TelegramAdapter{bot: bot, chatID: chatID}, nil
}

func (t *TelegramAdapter) SendMessage(ctx context.Context, message string) error {
	msg := tgbotapi.NewMessage(t.chatID, message)

	_, err := t.bot.Send(msg)

	return err
}

func (t *TelegramAdapter) SendMessageToChat(ctx context.Context, message string, chatID int64) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := t.bot.Send(msg)
	return err
}

func (t *TelegramAdapter) Start(ctx context.Context, messageHandler func(string, int64) error) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			message := strings.TrimSpace(update.Message.Text)
			if message == "" {
				continue
			}

			chatID := update.Message.Chat.ID
			if err := messageHandler(message, chatID); err != nil {
				// Отправляем сообщение об ошибке пользователю
				errorMsg := "Произошла ошибка при обработке запроса. Попробуйте еще раз."
				if err := t.SendMessageToChat(ctx, errorMsg, chatID); err != nil {
					// Логируем ошибку отправки сообщения об ошибке
					continue
				}
			}
		}
	}
}
