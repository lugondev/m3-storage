package user

type Service interface {
	// Authentication
	Authenticate(email string, password string) (*User, error)
	ValidateAPIKey(apiKey string) (*User, error)
	GetUser(id uint) (*User, error)

	// Quota management
	UpdateQuota(id uint, newQuota int64) error
	DeductQuota(id uint, amount int64) error
}
