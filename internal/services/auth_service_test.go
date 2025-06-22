package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hinphansa/7-solutions-challenge/internal/domain"
	"github.com/hinphansa/7-solutions-challenge/internal/mocks"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestAuthService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordHasher := mocks.NewMockPasswordHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	authService := NewAuthService(userRepo, passwordHasher, tokenGenerator)

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
	passwordHasher.EXPECT().Compare(gomock.Eq(user.Password), gomock.Eq(user.Password)).Return(nil).AnyTimes()
	tokenGenerator.EXPECT().Generate(gomock.Eq(user.ID), gomock.Eq(user.Email)).Return("token", nil).AnyTimes()

	token, err := authService.Login(context.Background(), user.Email, user.Password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token != "token" {
		t.Fatalf("expected token %v, got %v", "token", token)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordHasher := mocks.NewMockPasswordHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	authService := NewAuthService(userRepo, passwordHasher, tokenGenerator)

	userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq("test@example.com")).Return(nil, errors.New("user not found"))

	_, err := authService.Login(context.Background(), "test@example.com", "password")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != errUserNotFound.Error() {
		t.Fatalf("expected error %v, got %v", "user not found", err)
	}
}

func TestAuthService_Login_PasswordCompareError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordHasher := mocks.NewMockPasswordHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	authService := NewAuthService(userRepo, passwordHasher, tokenGenerator)

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
	passwordHasher.EXPECT().Compare(gomock.Eq(user.Password), gomock.Eq(user.Password)).Return(errors.New("password compare error")).AnyTimes()

	_, err := authService.Login(context.Background(), user.Email, user.Password)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != errInvalidPassword.Error() {
		t.Fatalf("expected error %v, got %v", "password compare error", err)
	}
}

func TestAuthService_Login_TokenGeneratorError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordHasher := mocks.NewMockPasswordHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	authService := NewAuthService(userRepo, passwordHasher, tokenGenerator)

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
	passwordHasher.EXPECT().Compare(gomock.Eq(user.Password), gomock.Eq(user.Password)).Return(nil).AnyTimes()
	tokenGenerator.EXPECT().Generate(gomock.Eq(user.ID), gomock.Eq(user.Email)).Return("", errors.New("token generator error")).AnyTimes()

	_, err := authService.Login(context.Background(), user.Email, user.Password)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != errUnableToGenerateToken.Error() {
		t.Fatalf("expected error %v, got %v", "token generator error", err)
	}
}
