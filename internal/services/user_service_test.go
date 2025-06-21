package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hinphansa/7-solutions-challenge/internal/domain"
	"github.com/hinphansa/7-solutions-challenge/internal/mocks"
	"github.com/hinphansa/7-solutions-challenge/internal/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestUserService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordHasher := mocks.NewMockPasswordHasher(ctrl)
	userService := NewUserService(userRepo, passwordHasher, nil) // tokenGenerator is nil because we don't need it for this test

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	id := bson.NewObjectID()
	userRepo.EXPECT().Create(gomock.Any(), gomock.Eq(user)).Return(&id, nil).AnyTimes()

	hashedPassword := "hashed_password"
	passwordHasher.EXPECT().Hash(user.Password).Return(hashedPassword, nil).AnyTimes()

	newID, err := userService.Register(context.Background(), user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newID == nil {
		t.Fatalf("expected id, got nil")
	}

	if newID.Hex() != id.Hex() {
		t.Fatalf("expected id %v, got %v", id, newID)
	}
}

func TestUserService_Register_Hash_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordHasher := mocks.NewMockPasswordHasher(ctrl)
	userService := NewUserService(userRepo, passwordHasher, nil) // tokenGenerator is nil because we don't need it for this test

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	passwordHasher.EXPECT().Hash(user.Password).Return("", errors.New("hash error"))

	_, err := userService.Register(context.Background(), user)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "hash error" {
		t.Fatalf("expected error %v, got %v", "hash error", err)
	}
}

func TestUserService_Register_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordHasher := mocks.NewMockPasswordHasher(ctrl)
	userService := NewUserService(userRepo, passwordHasher, nil) // tokenGenerator is nil because we don't need it for this test

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	hashedPassword := "hashed_password"
	passwordHasher.EXPECT().Hash(user.Password).Return(hashedPassword, nil).AnyTimes()

	userRepo.EXPECT().Create(gomock.Any(), gomock.Eq(user)).Return(nil, errors.New("create error"))

	_, err := userService.Register(context.Background(), user)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "create error" {
		t.Fatalf("expected error %v, got %v", "create error", err)
	}
}

func TestUserService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordHasher := mocks.NewMockPasswordHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	userService := NewUserService(userRepo, passwordHasher, tokenGenerator)

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
	passwordHasher.EXPECT().Compare(gomock.Eq(user.Password), gomock.Eq(user.Password)).Return(nil).AnyTimes()
	tokenGenerator.EXPECT().Generate(gomock.Eq(user.ID), gomock.Eq(user.Email)).Return("token", nil).AnyTimes()

	token, err := userService.Login(context.Background(), user.Email, user.Password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token != "token" {
		t.Fatalf("expected token %v, got %v", "token", token)
	}
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordHasher := mocks.NewMockPasswordHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	userService := NewUserService(userRepo, passwordHasher, tokenGenerator)

	userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq("test@example.com")).Return(nil, errors.New("user not found"))

	_, err := userService.Login(context.Background(), "test@example.com", "password")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != errUserNotFound.Error() {
		t.Fatalf("expected error %v, got %v", "user not found", err)
	}
}

func TestUserService_Login_PasswordCompareError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordHasher := mocks.NewMockPasswordHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	userService := NewUserService(userRepo, passwordHasher, tokenGenerator)

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
	passwordHasher.EXPECT().Compare(gomock.Eq(user.Password), gomock.Eq(user.Password)).Return(errors.New("password compare error")).AnyTimes()

	_, err := userService.Login(context.Background(), user.Email, user.Password)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != errInvalidPassword.Error() {
		t.Fatalf("expected error %v, got %v", "password compare error", err)
	}
}

func TestUserService_Login_TokenGeneratorError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordHasher := mocks.NewMockPasswordHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	userService := NewUserService(userRepo, passwordHasher, tokenGenerator)

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	userRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Eq(user.Email)).Return(user, nil).AnyTimes()
	passwordHasher.EXPECT().Compare(gomock.Eq(user.Password), gomock.Eq(user.Password)).Return(nil).AnyTimes()
	tokenGenerator.EXPECT().Generate(gomock.Eq(user.ID), gomock.Eq(user.Email)).Return("", errors.New("token generator error")).AnyTimes()

	_, err := userService.Login(context.Background(), user.Email, user.Password)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != errUnableToGenerateToken.Error() {
		t.Fatalf("expected error %v, got %v", "token generator error", err)
	}
}

func TestUserService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo, nil, nil) // passwordHasher and tokenGenerator are nil because we don't need them for this test

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	userRepo.EXPECT().GetByID(gomock.Any(), gomock.Eq(user.ID)).Return(user, nil).AnyTimes()

	foundUser, err := userService.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if foundUser == nil {
		t.Fatalf("expected user, got nil")
	}

	if foundUser.ID != user.ID {
		t.Fatalf("expected user %v, got %v", user, user)
	}
}

func TestUserService_GetByID_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo, nil, nil) // passwordHasher and tokenGenerator are nil because we don't need them for this test

	userRepo.EXPECT().GetByID(gomock.Any(), gomock.Eq(bson.ObjectID{})).Return(nil, errors.New("user not found"))

	_, err := userService.GetByID(context.Background(), bson.ObjectID{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != errUserNotFound.Error() {
		t.Fatalf("expected error %v, got %v", "user not found", err)
	}
}

func TestUserService_GetAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo, nil, nil) // passwordHasher and tokenGenerator are nil because we don't need them for this test

	userRepo.EXPECT().GetAll(gomock.Any()).Return([]domain.User{
		{
			ID:       bson.ObjectID{},
			Email:    "test@example.com",
			Password: "password",
		},
		{
			ID:       bson.ObjectID{},
			Email:    "test2@example.com",
			Password: "password2",
		},
	}, nil).AnyTimes()

	users, err := userService.GetAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %v", len(users))
	}
}

func TestUserService_GetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo, nil, nil) // passwordHasher and tokenGenerator are nil because we don't need them for this test

	userRepo.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("get all error"))

	_, err := userService.GetAll(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "get all error" {
		t.Fatalf("expected error %v, got %v", "get all error", err)
	}
}

func TestUserService_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo, nil, nil) // passwordHasher and tokenGenerator are nil because we don't need them for this test

	userRepo.EXPECT().List(gomock.Any(), gomock.Eq(&ports.Pagination{})).Return([]domain.User{
		{
			ID:       bson.ObjectID{},
			Email:    "test@example.com",
			Password: "password",
		},
	}, nil).AnyTimes()

	users, err := userService.List(context.Background(), &ports.Pagination{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %v", len(users))
	}
}

func TestUserService_List_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo, nil, nil) // passwordHasher and tokenGenerator are nil because we don't need them for this test

	userRepo.EXPECT().List(gomock.Any(), gomock.Eq(&ports.Pagination{})).Return(nil, errors.New("list error"))

	_, err := userService.List(context.Background(), &ports.Pagination{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "list error" {
		t.Fatalf("expected error %v, got %v", "list error", err)
	}
}

func TestUserService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo, nil, nil) // passwordHasher and tokenGenerator are nil because we don't need them for this test

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	userRepo.EXPECT().Update(gomock.Any(), gomock.Eq(user.ID), gomock.Eq(user)).Return(nil).AnyTimes()

	err := userService.Update(context.Background(), user.ID, user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserService_Update_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo, nil, nil) // passwordHasher and tokenGenerator are nil because we don't need them for this test

	user := &domain.User{
		ID:       bson.ObjectID{},
		Email:    "test@example.com",
		Password: "password",
	}

	userRepo.EXPECT().Update(gomock.Any(), gomock.Eq(user.ID), gomock.Eq(user)).Return(errors.New("update error")).AnyTimes()

	err := userService.Update(context.Background(), user.ID, user)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "update error" {
		t.Fatalf("expected error %v, got %v", "update error", err)
	}
}

func TestUserService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo, nil, nil) // passwordHasher and tokenGenerator are nil because we don't need them for this test

	userRepo.EXPECT().Delete(gomock.Any(), gomock.Eq(bson.ObjectID{})).Return(nil).AnyTimes()

	err := userService.Delete(context.Background(), bson.ObjectID{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserService_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepo, nil, nil) // passwordHasher and tokenGenerator are nil because we don't need them for this test

	userRepo.EXPECT().Delete(gomock.Any(), gomock.Eq(bson.ObjectID{})).Return(errors.New("delete error")).AnyTimes()

	err := userService.Delete(context.Background(), bson.ObjectID{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "delete error" {
		t.Fatalf("expected error %v, got %v", "delete error", err)
	}
}
