package adapter

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/modules/notify/port"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramAdapter implements the port.NotifyAdapter
type TelegramAdapter struct {
	bot         *tgbotapi.BotAPI
	chatID      int64
	logger      logger.Logger
	serviceName string
}

// TelegramConfig holds configuration for the Telegram adapter.
// TODO: Consider adding an 'Enabled' flag here for explicit control.
type TelegramConfig struct {
	BotToken string
	ChatID   string
	Debug    bool // Optional: Enable debug mode for the bot API
}

// NewTelegramAdapter creates a new Telegram adapter.
// Returns both port.NotifyAdapter and ports.HealthChecker.
// Consider splitting initialization if only one interface is needed sometimes.
func NewTelegramAdapter(config TelegramConfig, logger logger.Logger) (port.NotifyAdapter, error) {
	if config.BotToken == "" {
		return nil, fmt.Errorf("telegram bot token is required")
	}
	namedLogger := logger.WithFields(map[string]any{
		"service": "telegram_notify_adapter",
	})
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		namedLogger.Error(context.Background(), "Error creating Telegram bot API instance", map[string]any{"error": err})
		return nil, fmt.Errorf("failed to initialize Telegram bot: %w", err)
	}
	bot.Debug = config.Debug
	adapter := &TelegramAdapter{
		bot:         bot,
		logger:      namedLogger,
		serviceName: "telegram",
	}
	// Check if default chat ID is set
	if config.ChatID != "" {
		parsedChatID, err := strconv.ParseInt(config.ChatID, 10, 64)
		if err != nil {
			namedLogger.Error(context.Background(), "Invalid ChatID in Telegram config", map[string]any{"chat_id": config.ChatID, "error": err})
		} else {
			adapter.chatID = parsedChatID
			namedLogger.Info(context.Background(), "Telegram adapter configured with default chat ID", map[string]any{"chat_id": parsedChatID})
		}
	}

	namedLogger.Info(context.Background(), "Telegram notify adapter (library) initialized", map[string]any{"bot_username": bot.Self.UserName})
	return adapter, nil
}

// Send a message via Telegram using the library.
func (a *TelegramAdapter) Send(ctx context.Context, notification port.Content) error {

	// Message content can use Subject and Body combined, or just Body
	// For simplicity, using Body for now. Could be enhanced based on notification structure.
	messageText := notification.Message
	if notification.Subject != "" {
		messageText = fmt.Sprintf("Subject: %s\n\n%s", notification.Subject, notification.Message)
	}

	// Create the message object
	msg := tgbotapi.NewMessage(a.chatID, messageText)
	// Consider adding ParseMode (HTML or MarkdownV2) based on notification needs/data
	msg.ParseMode = tgbotapi.ModeHTML // Example

	// Send the message
	if _, sendErr := a.bot.Send(msg); sendErr != nil {
		a.logger.Error(ctx, "Error sending message via Telegram library", map[string]any{"chat_id": a.chatID, "error": sendErr})
	} else {
		a.logger.Debug(ctx, "Message sent successfully via Telegram library", map[string]any{"chat_id": a.chatID})
	}
	return nil
}
