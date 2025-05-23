package adapter

import (
	"context"

	"github.com/lugondev/go-log"
	"github.com/lugondev/m3-storage/internal/modules/notify/port"
)

// MockLogAdapter implements the port.NotifyAdapter interface by logging notifications.
type MockLogAdapter struct {
	logger logger.Logger
}

// NewMockLogAdapter creates a new instance of MockLogAdapter.
func NewMockLogAdapter(logger logger.Logger) port.NotifyAdapter {
	namedLogger := logger.WithFields(map[string]any{
		"service": "mocklog_notify_adapter",
	})
	ctx := context.Background()
	namedLogger.Info(ctx, "MockLog notify adapter initialized")

	return &MockLogAdapter{logger: namedLogger}
}

// Send logs the notification details.
func (a *MockLogAdapter) Send(ctx context.Context, notification port.Content) error {
	a.logger.Info(ctx, "--- MOCK Notification Sent (via Log) ---", map[string]any{
		"subject": notification.Subject,
		"message": notification.Message,
	})
	return nil
}
