package user

type Repository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByAPIKey(apiKey string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List() ([]*User, error)
}
