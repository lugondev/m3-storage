package domain

import (
	"time"

	"github.com/google/uuid"
)

// Media represents the metadata for an uploaded file.
type Media struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	FileName   string    `json:"file_name" gorm:"type:varchar(255)"`
	FilePath   string    `json:"file_path" gorm:"type:varchar(500)"` // Path in the adapters provider
	FileSize   int64     `json:"file_size"`
	MediaType  string    `json:"media_type" gorm:"type:varchar(50)"` // e.g., image, video, document
	Provider   string    `json:"provider" gorm:"type:varchar(50)"`   // e.g., local, s3, azure, firebase
	PublicURL  string    `json:"public_url" gorm:"type:varchar(500)"`
	UploadedAt time.Time `json:"uploaded_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName specifies the table name for the Media model.
func (Media) TableName() string {
	return "media"
}

// NewMedia creates a new Media entity.
func NewMedia(userID uuid.UUID, fileName, filePath string, fileSize int64, mediaType, provider, publicURL string) *Media {
	return &Media{
		ID:         uuid.New(),
		UserID:     userID,
		FileName:   fileName,
		FilePath:   filePath,
		FileSize:   fileSize,
		MediaType:  mediaType,
		Provider:   provider,
		PublicURL:  publicURL,
		UploadedAt: time.Now(),
	}
}
