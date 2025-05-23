package service

import (
	"context"
	"fmt"

	logger "github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/infra/config"

	"github.com/lugondev/m3-storage/internal/modules/notify/adapter"
	"github.com/lugondev/m3-storage/internal/modules/notify/port"
)

// notifyService implements the NotifyService interface.
// It holds a map of registered adapters keyed by channel name.
type notifyService struct {
	adapter port.NotifyAdapter
	logger  logger.Logger
	name    config.NotifyChannel
	ctx     context.Context
}

// NewNotifyService creates a new instance of NotifyService.
func NewNotifyService(cfg *config.Config, logger logger.Logger) (port.NotifyService, error) {
	namedLogger := logger.WithFields(map[string]any{
		"service": "notify_service_" + cfg.Adapter.Notify,
	})

	ctx := context.Background()
	logger.Debug(ctx, "Registered notify adapter", map[string]any{
		"channel": cfg.Adapter.Notify,
	})

	var notifyAdapter port.NotifyAdapter
	if cfg.Adapter.Notify == config.NotifyTelegram {
		telegramAdapter, err := adapter.NewTelegramAdapter(adapter.TelegramConfig{
			BotToken: cfg.Telegram.BotToken,
			ChatID:   cfg.Telegram.ChatID,
			Debug:    cfg.Telegram.Debug,
		}, logger)
		if err != nil {
			namedLogger.Error(ctx, "Failed to create Telegram adapter", map[string]any{
				"error": err,
			})
			return nil, fmt.Errorf("failed to create Telegram adapter: %w", err)
		}

		notifyAdapter = telegramAdapter
		namedLogger.Info(ctx, "Using Telegram adapter for notifications", map[string]any{
			"chat_id": cfg.Telegram.ChatID,
		})
	} else {
		notifyAdapter = adapter.NewMockLogAdapter(namedLogger)
		namedLogger.Info(ctx, "Using MockLog adapter for notifications", map[string]any{
			"channel": cfg.Adapter.Notify,
		})
	}

	return &notifyService{
		adapter: notifyAdapter,
		logger:  namedLogger,
		ctx:     ctx,
		name:    cfg.Adapter.Notify,
	}, nil
}

// Send finds the appropriate adapter based on the notification's channel
func (s *notifyService) Send(ctx context.Context, content port.Content) error {
	if content.Message == "" {
		return fmt.Errorf("notification message cannot be empty")
	}

	s.logger.Info(s.ctx, "Sending notification via adapter", map[string]any{
		"sub": content.Subject, // Subject might be empty
		"msg": content.Message,
	})

	err := s.adapter.Send(ctx, content)
	if err != nil {
		s.logger.Error(s.ctx, "Failed to send notification", map[string]any{
			"error": err,
		})
		return fmt.Errorf("failed to send notification: %w", err)
	}

	s.logger.Info(s.ctx, "Notification sent successfully")
	return nil
}

// ServiceName returns the name of the notification service.
func (s *notifyService) ServiceName() string {
	return string(s.name)
}
