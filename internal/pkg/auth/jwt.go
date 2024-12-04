package auth

import (
	"fmt"
	"strings"
	"time"

	"go-api-mono/internal/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

// Config JWT配置
type Config struct {
	SigningKey      string        // 签名密钥
	ExpirationTime  time.Duration // 过期时间
	SigningMethod   string        // 签名方法
	TokenPrefix     string        // 令牌前缀（如 "Bearer"）
	ClaimsKey       string        // Claims在上下文中的键
	ContextUserKey  string        // 用户信息在上下文中的键
	ContextTokenKey string        // 令牌在上下文中的键
}

// Claims 自定义JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`  // 用户ID
	Username string `json:"username"` // 用户名
	Role     string `json:"role"`     // 用户角色
	jwt.RegisteredClaims
}

// JWT 认证器
type JWT struct {
	Config Config
}

// New 创建一个新的JWT认证器
func New(config Config) *JWT {
	// 设置默认值
	if config.SigningMethod == "" {
		config.SigningMethod = "HS256"
	}
	if config.TokenPrefix == "" {
		config.TokenPrefix = "Bearer"
	}
	if config.ExpirationTime == 0 {
		config.ExpirationTime = 24 * time.Hour
	}
	if config.ClaimsKey == "" {
		config.ClaimsKey = "claims"
	}
	if config.ContextUserKey == "" {
		config.ContextUserKey = "user"
	}
	if config.ContextTokenKey == "" {
		config.ContextTokenKey = "token"
	}
	return &JWT{Config: config}
}

// GenerateToken 生成JWT令牌
func (j *JWT) GenerateToken(userID uint, username, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.Config.ExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(j.Config.SigningMethod), claims)
	signedToken, err := token.SignedString([]byte(j.Config.SigningKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ParseToken 解析JWT令牌
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	// 移除令牌前缀
	tokenString = j.stripTokenPrefix(tokenString)

	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != j.Config.SigningMethod {
			return nil, errors.ErrInvalidSignature
		}
		return []byte(j.Config.SigningKey), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, errors.ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, errors.New(errors.ErrCodeInvalidToken, "token not valid yet")
		default:
			return nil, errors.ErrInvalidToken
		}
	}

	if !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	return &claims, nil
}

// RefreshToken 刷新JWT令牌
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		// 只允许刷新过期的令牌
		if !errors.Is(err, errors.ErrTokenExpired) {
			return "", fmt.Errorf("failed to parse token for refresh: %w", err)
		}
	}

	// 更新时间声明
	now := time.Now()
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(j.Config.ExpirationTime))
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)

	// 生成新令牌
	token := jwt.NewWithClaims(jwt.GetSigningMethod(j.Config.SigningMethod), claims)
	signedToken, err := token.SignedString([]byte(j.Config.SigningKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign refreshed token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken 验证JWT令牌
func (j *JWT) ValidateToken(tokenString string) error {
	_, err := j.ParseToken(tokenString)
	return err
}

// GetUserFromToken 从令牌中获取用户信息
func (j *JWT) GetUserFromToken(tokenString string) (*Claims, error) {
	return j.ParseToken(tokenString)
}

// stripTokenPrefix 移除令牌前缀
func (j *JWT) stripTokenPrefix(tokenString string) string {
	if j.Config.TokenPrefix != "" && strings.HasPrefix(tokenString, j.Config.TokenPrefix+" ") {
		return strings.TrimPrefix(tokenString, j.Config.TokenPrefix+" ")
	}
	return tokenString
}
