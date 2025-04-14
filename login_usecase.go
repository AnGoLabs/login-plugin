// login-plugin/login_usecase.go
package login

import (
	"errors"
)

type TokenGenerator interface {
	GenerateAccessToken(userID, email, accountType string) (string, error)
	GenerateRefreshToken(userID string) (*RefreshToken, error)
}

type LoginUseCase struct {
	UserRepo      UserRepository
	RefreshRepo   RefreshTokenRepository
	TokenService  TokenGenerator
	PasswordCheck func(password, hash string) bool
}

func (uc *LoginUseCase) Execute(input LoginInput) (*LoginOutput, error) {
	user, err := uc.UserRepo.FindByEmail(input.Email)
	if err != nil || user == nil {
		return nil, errors.New("e-mail ou senha inválidos")
	}

	if !uc.PasswordCheck(input.Password, user.Password) {
		return nil, errors.New("e-mail ou senha inválidos")
	}

	accessToken, err := uc.TokenService.GenerateAccessToken(user.ID, user.Email, user.AccountType)
	if err != nil {
		return nil, errors.New("erro ao gerar access token")
	}

	refreshToken, err := uc.TokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("erro ao gerar refresh token")
	}

	_ = uc.RefreshRepo.DeleteByUserID(user.ID)
	err = uc.RefreshRepo.Save(refreshToken)
	if err != nil {
		return nil, errors.New("erro ao salvar refresh token")
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		User:         user,
	}, nil
}
