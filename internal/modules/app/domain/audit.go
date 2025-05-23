package domain

import (
	"time"

	"github.com/google/uuid"
)

// ActionType represents the type of action performed
type ActionType string

const (
	ActionTypeCreate ActionType = "create"
	ActionTypeRead   ActionType = "read"
	ActionTypeUpdate ActionType = "update"
	ActionTypeDelete ActionType = "delete"
	ActionTypeLogin  ActionType = "login"
	ActionTypeLogout ActionType = "logout"
)

// ResourceType represents the type of resource being acted upon
type ResourceType string

const (
	// System resource types
	ResourceTypeTenant ResourceType = "tenants"
	ResourceTypeUser   ResourceType = "users"
	ResourceTypeRole   ResourceType = "roles"

	// Tenant-specific resource types
	ResourceTypeTenantUser   ResourceType = "tenant_users"
	ResourceTypeTenantRole   ResourceType = "tenant_roles"
	ResourceTypeTenantConfig ResourceType = "tenant_config"

	// User-specific resource types
	ResourceTypeUserAuthentication ResourceType = "user_authentication"
	ResourceTypeUserProfile        ResourceType = "user_profile"
	ResourceTypeUserSettings       ResourceType = "user_settings"
	ResourceTypeUserActivity       ResourceType = "user_activity"
	ResourceTypeUserPreferences    ResourceType = "user_preferences"

	// Other resource types
	ResourceTypeAPIKey       ResourceType = "api_key"
	ResourceTypeAuditLog     ResourceType = "audit_log"
	ResourceTypeNotification ResourceType = "notification"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           uuid.UUID    `json:"id"`
	UserID       uuid.UUID    `json:"user_id"`             // User who performed the action
	TenantID     *uuid.UUID   `json:"tenant_id,omitempty"` // Optional: Tenant context of the action
	ActionType   ActionType   `json:"action_type"`
	ResourceType ResourceType `json:"resource_type"`
	ResourceID   string       `json:"resource_id"`
	Description  string       `json:"description"`
	Metadata     string       `json:"metadata,omitempty"` // JSON string with additional data
	IPAddress    string       `json:"ip_address"`
	UserAgent    string       `json:"user_agent"`
	CreatedAt    time.Time    `json:"created_at"`
}
