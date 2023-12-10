package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"

	"github.com/olad5/AfriHacks2023-stressless-backend/config"
	authMiddleware "github.com/olad5/AfriHacks2023-stressless-backend/internal/handlers/auth"
	loggingMiddleware "github.com/olad5/AfriHacks2023-stressless-backend/internal/handlers/logging"
	userHandlers "github.com/olad5/AfriHacks2023-stressless-backend/internal/handlers/users"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/infra/mongo"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/infra/redis"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/services/auth"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/usecases/users"
	"github.com/olad5/AfriHacks2023-stressless-backend/pkg/utils/logger"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

func NewHttpRouter(ctx context.Context, configurations *config.Configurations, logger *zap.Logger) http.Handler {
	opts := options.Client()
	mongoClient, err := mongoDriver.Connect(ctx, opts.ApplyURI(configurations.DatabaseUrl))
	if err != nil {
		log.Fatal("failed to create a mongo client: %w", err)
	}
	mongoDatabase := mongoClient.Database(configurations.DatabaseName)

	userRepo, err := mongo.NewMongoUserRepo(ctx, mongoDatabase, logger)
	if err != nil {
		log.Fatal("Error Initializing User Repo", err)
	}

	redisCache, err := redis.New(ctx, configurations, logger)
	if err != nil {
		log.Fatal("Error Initializing redisCache", err)
	}

	authService, err := auth.NewRedisAuthService(ctx, redisCache, configurations, logger)
	if err != nil {
		log.Fatal("Error Initializing Auth Service", err)
	}

	metricRepo, err := mongo.NewMongoMetricRepo(ctx, mongoDatabase, logger)
	if err != nil {
		log.Fatal("Error Initializing Metric Repo", err)
	}

	userService, err := users.NewUserService(userRepo, authService, metricRepo, logger)
	if err != nil {
		log.Fatal("Error Initializing UserService")
	}

	userHandler, err := userHandlers.NewUserHandler(*userService, authService, logger)
	if err != nil {
		log.Fatal("failed to create the User handler: ", err)
	}

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "StressLess Backend is live!")
	})

	router.Group(func(r chi.Router) {
		r.Use(
			middleware.AllowContentType("application/json"),
			middleware.SetHeader("Content-Type", "application/json"),
		)
		r.Post("/users/login", userHandler.Login)
		r.Post("/users", userHandler.CreateUser)
	})

	// -------------------------------------------------------------------------

	router.Group(func(r chi.Router) {
		r.Use(
			middleware.AllowContentType("application/json"),
			middleware.SetHeader("Content-Type", "application/json"),
		)
		r.Use(authMiddleware.EnsureAuthenticated(authService))

		r.Get("/users/me", userHandler.GetLoggedInUser)
		r.Patch("/users/onboarding", userHandler.CompleteOnboarding)
	})

	return router
}

func main() {
	configurations := config.GetConfig(".env")
	ctx := context.Background()

	l := logger.Get(configurations)
	appRouter := NewHttpRouter(ctx, configurations, l)

	port := configurations.Port

	server := &http.Server{Addr: ":" + port, Handler: loggingMiddleware.RequestLogger(appRouter, configurations)}
	go func() {
		l.Info(
			"starting application server on port: "+port,
			zap.String("port", port),
		)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exiting gracefully")
}
