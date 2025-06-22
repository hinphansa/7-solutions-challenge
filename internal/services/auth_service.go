package services

import (
	"context"

	"github.com/hinphansa/7-solutions-challenge/internal/ports"
)

var _ ports.AuthService = &authsvc{}

type authsvc struct {
	userRepo       ports.UserRepository
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
}

func NewAuthService(userRepo ports.UserRepository, passwordHasher PasswordHasher, tokenGenerator TokenGenerator) *authsvc {
	return &authsvc{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

func (s *authsvc) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errUserNotFound
	}

	if err := s.passwordHasher.Compare(password, user.Password); err != nil {
		return "", errInvalidPassword
	}

	token, err := s.tokenGenerator.Generate(user.ID, user.Email)
	if err != nil {
		return "", errUnableToGenerateToken
	}

	return token, nil
}
