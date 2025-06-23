package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hinphansa/7-solutions-challenge/api/gen/user/github.com/hinphansa/7-solutions-challenge/api/gen/user"
	"github.com/hinphansa/7-solutions-challenge/config"
	"github.com/hinphansa/7-solutions-challenge/internal/adapters/auth"
	grpc_adapter "github.com/hinphansa/7-solutions-challenge/internal/adapters/grpc"
	mongo_repo "github.com/hinphansa/7-solutions-challenge/internal/adapters/mongo"
	"github.com/hinphansa/7-solutions-challenge/internal/ports"
	"github.com/hinphansa/7-solutions-challenge/internal/services"
	"github.com/hinphansa/7-solutions-challenge/pkg/logger"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	l := logger.New(logrus.DebugLevel).WithFields(logrus.Fields{
		"service": "grpc",
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		l.Fatalf("Failed to load config: %v", err)
	}

	/* -------------------------------- Mongo DB -------------------------------- */
	// mongo client
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		l.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	// test connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		l.Fatalf("Failed to ping MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)
	mongoDB := mongoClient.Database(cfg.Mongo.DB)

	/* -------------------------------- Cryptography ---------------------------- */
	// cryptography service
	passwordHasher := auth.NewBCrypt(cfg.PasswordHasher.Cost)
	tokenGenerator := auth.NewJWT(cfg.JWT.Secret, time.Duration(cfg.JWT.TTL)*time.Second)

	/* -------------------------------- User Service ---------------------------- */
	// user service
	userRepo := mongo_repo.NewUserRepository(mongoDB)
	userService := services.NewUserService(userRepo, passwordHasher, tokenGenerator)

	// auth service
	authService := services.NewAuthService(userRepo, passwordHasher, tokenGenerator)

	/* -------------------------------- gRPC Server ---------------------------- */
	// create a new gRPC server with auth interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_adapter.UnaryAuthInterceptor(tokenGenerator)),
	)

	// register reflection service
	reflection.Register(grpcServer)

	// register user service
	userServer := grpc_adapter.NewUserServer(l, userService, authService)
	user.RegisterUserServiceServer(grpcServer, userServer)

	// start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcServer.Port))
	if err != nil {
		l.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		l.Infof("Starting gRPC server on port %d", cfg.GrpcServer.Port)
		if err := grpcServer.Serve(lis); err != nil {
			l.Fatalf("Failed to serve: %v", err)
		}
	}()

	/* ------------------------------ Schedule Job ------------------------------ */
	// This could be works with a cron job as well.
	// I'm using a ticker since the job is not critical and low complexity.
	schedule(ctx, l, userService, 10*time.Second)

	/* ---------------------------- graceful shutdown --------------------------- */
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("Shutting down server...")
	time.Sleep(1 * time.Second) // just simulate some cleanup logic

	// cleanup server
	grpcServer.GracefulStop()

	// cleanup mongo client
	if err := mongoClient.Disconnect(ctx); err != nil {
		l.Errorf("Failed to disconnect from MongoDB: %v", err)
	}

	// cancel context, on this process, we could add some cleanup logic here if needed
	cancel()
}

// Spawn a goroutine that runs every 10 seconds and logs the number of users in the DB
func schedule(ctx context.Context, l logger.Logger, userService ports.UserService, period time.Duration) {
	ticker := time.NewTicker(period)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				l.Info("Stopping schedule")
				return
			case <-ticker.C:
				users, err := userService.GetAll(ctx)
				if err != nil {
					l.Errorf("Failed to get all users: %v", err)
				}
				l.Infof("Number of users in the DB: %d", len(users))
			}
		}
	}()
}
