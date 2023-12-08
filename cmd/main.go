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

	"github.com/go-chi/chi/v5"

	"github.com/olad5/AfriHacks2023-stressless-backend/config"
	authMiddleware "github.com/olad5/AfriHacks2023-stressless-backend/internal/handlers/auth"
	userHandlers "github.com/olad5/AfriHacks2023-stressless-backend/internal/handlers/users"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/infra/mongo"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/infra/redis"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/services/auth"
	"github.com/olad5/AfriHacks2023-stressless-backend/internal/usecases/users"
)

func NewHttpRouter(ctx context.Context, configurations *config.Configurations) http.Handler {
	userRepo, err := mongo.NewMongoUserRepo(ctx, configurations)
	if err != nil {
		log.Fatal("Error Initializing User Repo", err)
	}

	redisCache, err := redis.New(ctx, configurations)
	if err != nil {
		log.Fatal("Error Initializing redisCache", err)
	}

	authService, err := auth.NewRedisAuthService(ctx, redisCache, configurations)
	if err != nil {
		log.Fatal("Error Initializing Auth Service", err)
	}

	userService, err := users.NewUserService(userRepo, authService)
	if err != nil {
		log.Fatal("Error Initializing UserService")
	}

	userHandler, err := userHandlers.NewUserHandler(*userService, authService)
	if err != nil {
		log.Fatal("failed to create the User handler: ", err)
	}

	router := chi.NewRouter()

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
	})

	return router
}

func main() {
	configurations := config.GetConfig(".env")
	ctx := context.Background()

	appRouter := NewHttpRouter(ctx, configurations)

	port := configurations.Port
	server := &http.Server{Addr: ":" + port, Handler: appRouter}
	go func() {
		message := "Server is running on port " + port
		fmt.Println(message)
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
