// login-plugin/refresh_token_usecase.go
package login

import (
	"errors"
	"time"
)

type RefreshTokenUseCase struct {
	UserRepo     UserRepository
	RefreshRepo  RefreshTokenRepository
	TokenService TokenGenerator
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (uc *RefreshTokenUseCase) Execute(input RefreshInput) (*RefreshOutput, error) {
	token, err := uc.RefreshRepo.GetByToken(input.RefreshToken)
	if err != nil || token == nil {
		return nil, errors.New("refresh token inválido")
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expirado")
	}

	user, err := uc.UserRepo.FindByID(token.UserID)
	if err != nil || user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	access, err := uc.TokenService.GenerateAccessToken(user.ID, user.Email, user.AccountType)
	if err != nil {
		return nil, errors.New("erro ao gerar access token")
	}

	newRefresh, err := uc.TokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("erro ao gerar novo refresh token")
	}

	_ = uc.RefreshRepo.DeleteByUserID(user.ID)
	err = uc.RefreshRepo.Save(newRefresh)
	if err != nil {
		return nil, errors.New("erro ao salvar novo refresh token")
	}

	return &RefreshOutput{
		AccessToken:  access,
		RefreshToken: newRefresh.Token,
	}, nil
}
