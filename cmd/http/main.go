package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hinphansa/7-solutions-challenge/config"
	"github.com/hinphansa/7-solutions-challenge/internal/adapters/auth"
	"github.com/hinphansa/7-solutions-challenge/internal/adapters/http"
	mongo_repo "github.com/hinphansa/7-solutions-challenge/internal/adapters/mongo"
	"github.com/hinphansa/7-solutions-challenge/internal/ports"
	"github.com/hinphansa/7-solutions-challenge/internal/services"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/hinphansa/7-solutions-challenge/pkg/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	l := logger.New(logrus.DebugLevel).WithFields(logrus.Fields{
		"service": "http",
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

	// user handler
	userHandler := http.NewUserHandler(l, userService)

	/* -------------------------------- Auth Service ---------------------------- */
	// auth service
	authService := services.NewAuthService(userRepo, passwordHasher, tokenGenerator)

	// auth handler
	authHandler := http.NewAuthHandler(l, authService, userHandler)

	/* -------------------------------- Fiber app ------------------------------- */

	// create a new fiber app
	app := fiber.New()

	// apply general middlewares
	app.Use(http.RequestIdMiddleware())
	app.Use(http.LoggerMiddleware())

	// setup routes
	http.SetupRoutes(app, cfg, userHandler, authHandler)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.HttpServer.Port)); err != nil {
			l.Fatalf("Failed to start server: %v", err)
		}
	}()

	/* ------------------------------ Schedule Job ------------------------------ */
	// This could be works with a cron job as well.
	// I'm using a ticker since the job is not critical and low complexity.
	schedule(ctx, l, userService, 10*time.Second)

	/* ---------------------------- graceful shutdown --------------------------- */

	fmt.Println("Starting server...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("Shutting down server...")
	time.Sleep(1 * time.Second) // just simulate some cleanup logic

	// cleanup server
	if err := app.Shutdown(); err != nil {
		l.Errorf("Failed to shutdown server: %v", err)
	}

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
