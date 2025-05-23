package service

import (
	"context"
	"time"

	domains "github.com/lugondev/m3-storage/internal/modules/app/domain"
	ports "github.com/lugondev/m3-storage/internal/modules/app/port"

	"github.com/google/uuid"
)

type auditService struct {
	auditRepo ports.AuditRepository
}

func NewAuditService(auditRepo ports.AuditRepository) ports.AuditService {
	return &auditService{
		auditRepo: auditRepo,
	}
}

// Log creates a new audit log entry
func (s *auditService) Log(ctx context.Context, log *domains.AuditLog) error {
	// Set ID and timestamp if not already set
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}

	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}

	return s.auditRepo.Create(ctx, log)
}

// GetUserLogs retrieves audit logs for a specific user
func (s *auditService) GetUserLogs(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domains.AuditLog, error) {
	return s.auditRepo.GetByUser(ctx, userID, limit, offset)
}

// GetResourceLogs retrieves audit logs for a specific resource
func (s *auditService) GetResourceLogs(ctx context.Context, resourceType domains.ResourceType, resourceID string, limit, offset int) ([]domains.AuditLog, error) {
	return s.auditRepo.GetByResource(ctx, resourceType, resourceID, limit, offset)
}

// GetActionLogs retrieves audit logs for a specific action type
func (s *auditService) GetActionLogs(ctx context.Context, actionType domains.ActionType, limit, offset int) ([]domains.AuditLog, error) {
	return s.auditRepo.GetByAction(ctx, actionType, limit, offset)
}

// Search searches audit logs based on various criteria
func (s *auditService) Search(ctx context.Context, params map[string]any, limit, offset int) ([]domains.AuditLog, error) {
	return s.auditRepo.Search(ctx, params, limit, offset)
}
