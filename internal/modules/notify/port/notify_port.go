package port

import "context"

// Content represents the data structure for a generic notification.
// We might need more specific structures or fields depending on the channel.
type Content struct {
	Subject string
	Message string
}

// NotifyAdapter defines the interface for sending notifications via different providers/channels.
type NotifyAdapter interface {
	Send(ctx context.Context, content Content) error
}

// NotifyService defines the core logic for handling notifications.
type NotifyService interface {
	Send(ctx context.Context, content Content) error
	ServiceName() string
}
