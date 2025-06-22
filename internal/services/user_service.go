package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hinphansa/7-solutions-challenge/internal/domain"
	"github.com/hinphansa/7-solutions-challenge/internal/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Interface implementation check
var _ ports.UserService = &usersvc{}

var (
	errUserNotFound          = errors.New("user not found")
	errInvalidPassword       = errors.New("invalid password")
	errUnableToGenerateToken = errors.New("unable to generate token")
)

// PasswordHasher is an interface that defines the methods for hashing and comparing passwords
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) error
}

// TokenGenerator is an interface that defines the methods for generating and verifying tokens
type TokenGenerator interface {
	Generate(id bson.ObjectID, email string) (string, error)
}

type usersvc struct {
	userRepo       ports.UserRepository
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
}

func NewUserService(userRepo ports.UserRepository, passwordHasher PasswordHasher, tokenGenerator TokenGenerator) *usersvc {
	return &usersvc{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

func (s *usersvc) Register(ctx context.Context, user *domain.User) (*bson.ObjectID, error) {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.CreatedAt = time.Now()

	hash, err := s.passwordHasher.Hash(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hash

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s *usersvc) GetByID(ctx context.Context, id bson.ObjectID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errUserNotFound
	}
	return user, nil
}

func (s *usersvc) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *usersvc) List(ctx context.Context, pagination *ports.Pagination) ([]domain.User, error) {
	return s.userRepo.List(ctx, pagination)
}

func (s *usersvc) Update(ctx context.Context, id bson.ObjectID, user *domain.User) error {
	return s.userRepo.Update(ctx, id, user)
}

func (s *usersvc) Delete(ctx context.Context, id bson.ObjectID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *usersvc) Count(ctx context.Context) (int64, error) {
	return s.userRepo.Count(ctx)
}
