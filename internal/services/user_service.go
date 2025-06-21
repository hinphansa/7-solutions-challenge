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
var _ ports.UserService = &service{}

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

type service struct {
	userRepo       ports.UserRepository
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
}

func NewUserService(userRepo ports.UserRepository, passwordHasher PasswordHasher, tokenGenerator TokenGenerator) *service {
	return &service{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

func (s *service) Register(ctx context.Context, user *domain.User) (*bson.ObjectID, error) {
	// TODO: check email trimming and lowercasing, and validate email format
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

func (s *service) Login(ctx context.Context, email string, password string) (string, error) {
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

func (s *service) GetByID(ctx context.Context, id bson.ObjectID) (*domain.User, error) {
	// TODO: check if password is not returned
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errUserNotFound
	}
	return user, nil
}

func (s *service) GetAll(ctx context.Context) ([]domain.User, error) {
	// TODO: check if password is not returned
	return s.userRepo.GetAll(ctx)
}

func (s *service) List(ctx context.Context, pagination *ports.Pagination) ([]domain.User, error) {
	return s.userRepo.List(ctx, pagination)
}

func (s *service) Update(ctx context.Context, id bson.ObjectID, user *domain.User) error {
	// TODO: and authorization check
	// TODO: validate to allow only name or email to be updated
	return s.userRepo.Update(ctx, id, user)
}

func (s *service) Delete(ctx context.Context, id bson.ObjectID) error {
	// TODO: and authorization check
	return s.userRepo.Delete(ctx, id)
}
