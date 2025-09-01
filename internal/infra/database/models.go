package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lugondev/m3-storage/internal/shared/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StringArray represents a PostgreSQL string array
type StringArray []string

// Value converts StringArray to PostgreSQL string array
func (a StringArray) Value() (driver.Value, error) {
	return "{" + strings.Join(a, ",") + "}", nil
}

// Scan converts PostgreSQL string array to StringArray
func (a *StringArray) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		str := string(src)
		*a = strings.Split(strings.Trim(str, "{}"), ",")
		return nil
	case string:
		*a = strings.Split(strings.Trim(src, "{}"), ",")
		return nil
	case nil:
		*a = make(StringArray, 0)
		return nil
	}
	return fmt.Errorf("unsupported type for StringArray: %T", src)
}

// JSONB represents a PostgreSQL JSONB type
type JSONB json.RawMessage

// Value converts JSONB to database value
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// Scan converts database value to JSONB
func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.NewValidationError("invalid scan source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

// Base contains common columns for all tables
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// User represents a user in the system (globally)
type User struct {
	Base
	Email          string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash   string     `gorm:"type:text;not null"`
	FirstName      string     `gorm:"type:varchar(100);not null"`
	LastName       string     `gorm:"type:varchar(100);not null"`
	Status         string     `gorm:"type:varchar(20);not null;default:'active';index:idx_users_status"`
	EmailVerified  bool       `gorm:"not null;default:false"`
	LastLoginAt    *time.Time `gorm:"index:idx_users_last_login"`
	FailedAttempts int        `gorm:"not null;default:0"`
	LockedUntil    *time.Time
	Metadata       JSONB `gorm:"type:jsonb"`
}

// UserProfile represents additional user profile information
type UserProfile struct {
	Base
	UserID      uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Avatar      string    `gorm:"type:text"`
	PhoneNumber string    `gorm:"type:varchar(20)"`
	DateOfBirth *time.Time
	Timezone    string `gorm:"type:varchar(50);default:'UTC'"`
	Language    string `gorm:"type:varchar(10);default:'en'"`

	// Foreign key relationship
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// AuditLog represents system audit logs
type AuditLog struct {
	Base
	UserID       *uuid.UUID `gorm:"type:uuid;index:idx_audit_logs_user_id"`
	ActionType   string     `gorm:"type:varchar(50);not null;index:idx_audit_logs_action_type"`
	ResourceType string     `gorm:"type:varchar(50);not null;index:idx_audit_logs_resource_type"`
	ResourceID   string     `gorm:"type:varchar(255);not null"`
	Description  string     `gorm:"type:text"`
	Metadata     JSONB      `gorm:"type:jsonb"`
	IPAddress    string     `gorm:"type:varchar(45)"`
	UserAgent    string     `gorm:"type:text"`
}
