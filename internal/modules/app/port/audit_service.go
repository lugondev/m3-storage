package port

import (
	"context"

	domains "github.com/lugondev/m3-storage/internal/modules/app/domain"

	"github.com/google/uuid"
)

// AuditService defines the interface for audit logging
type AuditService interface {
	// Log creates a new audit log entry
	Log(ctx context.Context, log *domains.AuditLog) error

	// GetUserLogs retrieves audit logs for a specific user
	GetUserLogs(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domains.AuditLog, error)

	// GetResourceLogs retrieves audit logs for a specific resource
	GetResourceLogs(ctx context.Context, resourceType domains.ResourceType, resourceID string, limit, offset int) ([]domains.AuditLog, error)

	// GetActionLogs retrieves audit logs for a specific action type
	GetActionLogs(ctx context.Context, actionType domains.ActionType, limit, offset int) ([]domains.AuditLog, error)

	// Search searches audit logs based on various criteria
	Search(ctx context.Context, params map[string]any, limit, offset int) ([]domains.AuditLog, error)
}

// AuditRepository defines the interface for audit log adapters
type AuditRepository interface {
	// Create creates a new audit log entry
	Create(ctx context.Context, log *domains.AuditLog) error

	// GetByUser retrieves audit logs for a specific user
	GetByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domains.AuditLog, error)

	// GetByResource retrieves audit logs for a specific resource
	GetByResource(ctx context.Context, resourceType domains.ResourceType, resourceID string, limit, offset int) ([]domains.AuditLog, error)

	// GetByAction retrieves audit logs for a specific action type
	GetByAction(ctx context.Context, actionType domains.ActionType, limit, offset int) ([]domains.AuditLog, error)

	// Search searches audit logs based on various criteria
	Search(ctx context.Context, params map[string]any, limit, offset int) ([]domains.AuditLog, error)
}
