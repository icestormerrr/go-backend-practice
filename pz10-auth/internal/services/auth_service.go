package services

import (
	"errors"
	"time"

	"github.com/icestormerrr/pz10-auth/internal/core"
)

type AuthServiceConfig struct {
	RefreshTTL       time.Duration
	AccessTTL        time.Duration
	MaxLoginAttempts int64
}
type AuthService struct {
	config       AuthServiceConfig
	userRepo     core.UserRepo
	sessionRepo  core.SessionRepo
	tokenManager core.TokenManager
}

func NewAuthService(config AuthServiceConfig, u core.UserRepo, s core.SessionRepo, t core.TokenManager) *AuthService {
	return &AuthService{
		config:       config,
		userRepo:     u,
		sessionRepo:  s,
		tokenManager: t,
	}
}

func (s *AuthService) Login(email, password string) (accessToken, refreshToken string, userID int64, err error) {
	loginAttemptsCount, err := s.sessionRepo.IncLoginAttempts(email)
	if err != nil {
		return "", "", 0, errors.New("internal_error")
	}

	if loginAttemptsCount > s.config.MaxLoginAttempts {
		return "", "", 0, errors.New("too_many_login_attempts")
	}

	u, err := s.userRepo.CheckPassword(email, password)
	if err != nil {
		return "", "", 0, errors.New("unauthorized")
	}

	s.sessionRepo.ResetLoginAttempts(email)

	accessToken, err = s.tokenManager.Sign(u.ID, u.Email, u.Role, s.config.AccessTTL)
	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err = s.tokenManager.Sign(u.ID, u.Email, u.Role, s.config.RefreshTTL)
	if err != nil {
		return "", "", 0, err
	}

	if err := s.sessionRepo.SetRefreshToken(u.ID, refreshToken); err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, u.ID, nil
}

func (s *AuthService) RefreshTokens(oldRefreshToken string) (newAccessToken, newRefreshToken string, err error) {
	claims, err := s.tokenManager.Parse(oldRefreshToken)
	if err != nil {
		return "", "", errors.New("invalid_refresh_token")
	}

	userID := int64(claims["sub"].(float64))

	storedToken, err := s.sessionRepo.GetRefreshToken(userID)
	if err != nil || storedToken == "" {
		return "", "", errors.New("refresh_token_not_found")
	}

	if storedToken != oldRefreshToken {
		return "", "", errors.New("invalid_refresh_token")
	}

	u, err := s.userRepo.GetById(userID)
	if err != nil {
		return "", "", errors.New("user_not_found")
	}

	newAccessToken, err = s.tokenManager.Sign(u.ID, u.Email, u.Role, s.config.AccessTTL)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err = s.tokenManager.Sign(u.ID, u.Email, u.Role, s.config.RefreshTTL)
	if err != nil {
		return "", "", err
	}

	if err := s.sessionRepo.SetRefreshToken(userID, newRefreshToken); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
