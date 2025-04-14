# 🔐 AuthKit — Módulo de Autenticação Reutilizável para Go + Gin

> **Este projeto é de propriedade intelectual da AnGoLabs. Todos os direitos reservados.**

AuthKit é um plugin reutilizável de autenticação para projetos em Go usando Gin, JWT e GORM.

## ✨ Funcionalidades

- Login com JWT e Refresh Token
- RBAC (controle de acesso por tipo de conta)
- Middleware JWT + RoleRequired
- Expiração e renovação segura com rotação de refresh token
- Rate limit configurável para rota de login

---

## 📦 Instalação

```bash
go get github.com/seunome/login-plugin
```

> Certifique-se de ter iniciado um módulo Go com `go mod init` antes disso.

---

## ⚙️ Configuração básica

### 1. Implemente os contratos:

#### `UserRepository`
```go
type User struct {
  ID          string
  Email       string
  Password    string
  AccountType string
}

type UserRepository interface {
  FindByEmail(email string) (*User, error)
  FindByID(id string) (*User, error)
}
```

#### `RefreshTokenRepository`
```go
type RefreshToken struct {
  Token     string
  UserID    string
  ExpiresAt time.Time
}

type RefreshTokenRepository interface {
  Save(token *RefreshToken) error
  GetByToken(token string) (*RefreshToken, error)
  DeleteByUserID(userID string) error
}
```

### 2. Crie a instância

```go
authKit := login.NewAuthKit(login.Config{
  UserRepo:             seuUserRepo,
  RefreshTokenRepo:     seuRefreshTokenRepo,
  TokenSecret:          "sua-chave-secreta",
  AccessTokenDuration:  time.Hour,
  RefreshTokenDuration: 48 * time.Hour,
  PasswordCheckFunc:    func(p, h string) bool {
    return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
  },
})
```

### 3. Registre os endpoints

```go
authHandler := login.NewHandler(authKit.LoginUC, authKit.RefreshUC)
authMiddleware := login.NewAuthMiddleware("sua-chave-secreta")
rateLimiter := login.NewRateLimiter(1, 5, 10*time.Minute)

api := r.Group("/api")
{
  api.POST("/login", rateLimiter.LoginLimiterMiddleware(), authHandler.Login)
  api.POST("/refresh", authHandler.Refresh)
  api.GET("/me",
    authMiddleware.JWTAuth(),
    authMiddleware.RoleRequired("user", "company"),
    func(c *gin.Context) {
      c.JSON(200, gin.H{"confirm": "Tudo ok!"})
    },
  )
}
```

---

## 🔐 Segurança
- Tokens são validados com JWT HMAC-SHA256
- Refresh tokens são únicos e rotativos por usuário
- Rate limit por IP e e-mail previne ataques de força bruta

---

## 📌 Licença
MIT — uso autorizado apenas em projetos vinculados à AnGoLabs ou com autorização prévia.
