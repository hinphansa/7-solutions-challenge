package grpc

import (
	"context"

	"github.com/hinphansa/7-solutions-challenge/api/gen/user/github.com/hinphansa/7-solutions-challenge/api/gen/user"
	"github.com/hinphansa/7-solutions-challenge/internal/domain"
	"github.com/hinphansa/7-solutions-challenge/internal/ports"
	"github.com/hinphansa/7-solutions-challenge/pkg/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServer struct {
	user.UnimplementedUserServiceServer
	log         logger.Logger
	userService ports.UserService
	authService ports.AuthService
}

func NewUserServer(log logger.Logger, userService ports.UserService, authService ports.AuthService) *UserServer {
	return &UserServer{
		log:         log,
		userService: userService,
		authService: authService,
	}
}

// CreateUser implements the CreateUser RPC method
func (s *UserServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	id, err := s.userService.Register(ctx, &domain.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		s.log.Errorf("Failed to create user: %v", err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &user.CreateUserResponse{Id: id.Hex()}, nil
}

// GetUserById implements the GetUserById RPC method
func (s *UserServer) GetUserById(ctx context.Context, req *user.GetUserRequest) (*user.User, error) {
	id, err := bson.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	u, err := s.userService.GetByID(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to get user: %v", err)
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &user.User{
		Id:        u.ID.Hex(),
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: timestamppb.New(u.CreatedAt),
	}, nil
}

// ListUsers implements the ListUsers RPC method
func (s *UserServer) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	users, err := s.userService.List(ctx, &ports.Pagination{
		Limit:  int64(req.GetLimit()),
		Offset: int64(req.GetOffset()),
	})
	if err != nil {
		s.log.Errorf("Failed to list users: %v", err)
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	response := &user.ListUsersResponse{
		Users: make([]*user.User, len(users)),
	}

	for i, u := range users {
		response.Users[i] = &user.User{
			Id:        u.ID.Hex(),
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: timestamppb.New(u.CreatedAt),
		}
	}

	return response, nil
}

// Login implements the Login RPC method
func (s *UserServer) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	token, err := s.authService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		s.log.Errorf("Failed to login: %v", err)
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	return &user.LoginResponse{Token: token}, nil
}

// UpdateUser implements the UpdateUser RPC method
func (s *UserServer) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	// Get user ID from JWT token metadata
	userID, ok := ctx.Value(userIDKey).(bson.ObjectID)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing user ID")
	}

	reqID, err := bson.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	if len(userID) == 0 || userID != reqID {
		return nil, status.Error(codes.Unauthenticated, "missing user ID")
	}

	err = s.userService.Update(ctx, reqID, &domain.User{
		Name:  req.GetName(),
		Email: req.GetEmail(),
	})
	if err != nil {
		s.log.Errorf("Failed to update user: %v", err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return &user.UpdateUserResponse{
		Message: "user updated successfully",
	}, nil
}

// DeleteUser implements the DeleteUser RPC method
func (s *UserServer) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	// Get user ID from JWT token metadata
	userID, ok := ctx.Value(userIDKey).(bson.ObjectID)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing user ID")
	}

	reqID, err := bson.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	if len(userID) == 0 || userID != reqID {
		return nil, status.Error(codes.Unauthenticated, "missing user ID")
	}

	err = s.userService.Delete(ctx, reqID)
	if err != nil {
		s.log.Errorf("Failed to delete user: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	return &user.DeleteUserResponse{
		Message: "user deleted successfully",
	}, nil
}
