// login-plugin/auth_handler.go
package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	LoginUC   *LoginUseCase
	RefreshUC *RefreshTokenUseCase
}

func NewHandler(loginUC *LoginUseCase, refreshUC *RefreshTokenUseCase) *Handler {
	return &Handler{
		LoginUC:   loginUC,
		RefreshUC: refreshUC,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	output, err := h.LoginUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *Handler) Refresh(c *gin.Context) {
	var input RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	output, err := h.RefreshUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
