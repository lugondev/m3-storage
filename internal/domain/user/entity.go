package user

type User struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	ApiKey    string `json:"api_key"`
	Name      string `json:"name"`
	Quota     int64  `json:"quota"`
	Role      string `json:"role"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
