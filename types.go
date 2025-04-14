// login-plugin/types.go
package login

import "time"

type User struct {
	ID          string
	Email       string
	Password    string
	AccountType string
}

type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
}

type RefreshTokenRepository interface {
	Save(token *RefreshToken) error
	GetByToken(token string) (*RefreshToken, error)
	DeleteByUserID(userID string) error
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	User         *User
}
