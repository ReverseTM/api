package app

import (
	authHandler "api/internal/http/handlers/api/auth"
	managementHandler "api/internal/http/handlers/api/management"
	"api/internal/http/middlewares"
	"api/internal/http/server"
	"api/internal/services/auth"
	"api/internal/services/management"
	"api/internal/storage/mockUserStorage"
	"api/internal/storage/tarantool"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

type App struct {
	server  *server.Server
	storage *tarantool.Storage
}

func NewApp(
	log *slog.Logger,
	address string,
	storagePath string,
	spaceName string,
	secretKey string,
	tokenTTL time.Duration,
) *App {
	userStorage := mockUserStorage.NewUserStorage()
	storage, err := tarantool.NewTarantoolStorage(storagePath, spaceName)
	if err != nil {
		panic(err)
	}

	authService := auth.NewAuthService(log, userStorage, userStorage, secretKey, tokenTTL)
	managementService := management.NewManagementService(log, storage, storage)

	authH := authHandler.NewAuthHandler(log, authService)
	managementH := managementHandler.NewManagementHandler(log, managementService)

	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/login", authH.Login)
		api.POST("/read", middlewares.AuthMiddleware(secretKey), managementH.Read)
		api.POST("/write", middlewares.AuthMiddleware(secretKey), managementH.Write)
	}

	return &App{
		server:  server.NewServer(address, router),
		storage: storage,
	}
}

func (a *App) Run() {
	a.server.MustRun()
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.storage.Close(); err != nil {
		return err
	}

	return nil
}
