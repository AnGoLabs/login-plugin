// login-plugin/config.go
package login

import "time"

type Config struct {
	UserRepo             UserRepository
	RefreshTokenRepo     RefreshTokenRepository
	TokenSecret          string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	PasswordCheckFunc    func(password, hash string) bool
}

type AuthKit struct {
	LoginUC *LoginUseCase
}

func NewAuthKit(cfg Config) *AuthKit {
	tokenService := NewTokenService(TokenConfig{
		Secret:     cfg.TokenSecret,
		AccessTTL:  cfg.AccessTokenDuration,
		RefreshTTL: cfg.RefreshTokenDuration,
	})

	loginUC := &LoginUseCase{
		UserRepo:      cfg.UserRepo,
		RefreshRepo:   cfg.RefreshTokenRepo,
		TokenService:  tokenService,
		PasswordCheck: cfg.PasswordCheckFunc,
	}

	return &AuthKit{
		LoginUC: loginUC,
	}
}
