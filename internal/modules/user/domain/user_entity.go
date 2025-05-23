package domain

import (
	"time"
)

// User represents a user in the system.
type User struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"` // Store hashed password
	FirstName string    `json:"first_name" gorm:"type:varchar(100)"`
	LastName  string    `json:"last_name" gorm:"type:varchar(100)"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	APIKey    string    `json:"-" gorm:"type:varchar(255);uniqueIndex"` // Store hashed API Key or the key itself if managed externally
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Upload Quotas
	MaxStorageBytes    int64     `json:"max_storage_bytes" gorm:"default:1073741824"` // Default 1GB
	UsedStorageBytes   int64     `json:"used_storage_bytes" gorm:"default:0"`
	MaxFilesPerDay     int       `json:"max_files_per_day" gorm:"default:100"`
	UploadedFilesToday int       `json:"uploaded_files_today" gorm:"default:0"`
	LastUploadDate     time.Time `json:"last_upload_date"` // To reset UploadedFilesToday
}

// TableName specifies the table name for the User model.
func (User) TableName() string {
	return "users"
}
