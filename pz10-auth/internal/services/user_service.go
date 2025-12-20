package services

import (
	"errors"

	"github.com/icestormerrr/pz10-auth/internal/core"
)

type UserService struct {
	userRepo core.UserRepo
}

func NewUserService(u core.UserRepo) *UserService {
	return &UserService{userRepo: u}
}

func (s *UserService) GetById(userID int64) (core.User, error) {
	u, err := s.userRepo.GetById(userID)
	if err != nil {
		return core.User{}, errors.New("user_not_found")
	}
	return u, nil
}

func (s *UserService) GetStats() (map[string]any, error) {
	return map[string]any{"users": 2, "version": "1.0"}, nil
}
