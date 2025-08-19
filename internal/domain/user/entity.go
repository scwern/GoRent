package user

type Role string

const (
	AdminRole   Role = "admin"
	ManagerRole Role = "manager"
	ClientRole  Role = "client"
)

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         Role   `json:"role"`
}

func (u *User) IsAdmin() bool {
	return u.Role == AdminRole
}

func (u *User) IsManager() bool {
	return u.Role == ManagerRole
}

func (u *User) IsClient() bool {
	return u.Role == ClientRole
}
