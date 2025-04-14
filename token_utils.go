// login-plugin/token_utils.go
package login

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type TokenService struct {
	Config TokenConfig
}

func NewTokenService(config TokenConfig) *TokenService {
	return &TokenService{Config: config}
}

func (s *TokenService) GenerateAccessToken(userID, email, accountType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":      userID,
		"email":        email,
		"account_type": accountType,
		"exp":          time.Now().Add(s.Config.AccessTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Config.Secret))
}

func (s *TokenService) GenerateRefreshToken(userID string) (*RefreshToken, error) {
	return &RefreshToken{
		Token:     uuid.NewString(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.Config.RefreshTTL),
	}, nil
}
