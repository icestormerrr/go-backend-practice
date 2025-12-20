package core

import "time"

type UserService interface {
	GetById(userID int64) (User, error)
	GetStats() (map[string]any, error)
}

type AuthService interface {
	Login(email, password string) (accessToken, refreshToken string, userID int64, err error)
	RefreshTokens(oldRefreshToken string) (newAccessToken, newRefreshToken string, err error)
}

type TokenManager interface {
	Sign(userID int64, email, role string, ttl time.Duration) (string, error)
	Parse(tokenStr string) (map[string]any, error)
}
