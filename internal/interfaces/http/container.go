package http

import (
	"dashboard-starter/db"
	"dashboard-starter/internal/application/services"
	"dashboard-starter/internal/infrastructure"
	"dashboard-starter/internal/interfaces/http/handlers"
)

type Container struct {
	infraFactory *infrastructure.InfrastructureFactory

	// Services
	authService    *services.AuthService
	userService    *services.UserApplicationService
	deviceService  *services.DeviceApplicationService
	articleService *services.ArticleApplicationService

	// Handlers
	AuthHandler    *handlers.AuthHandler
	UserHandler    *handlers.UserHandler
	DeviceHandler  *handlers.DeviceHandler
	ArticleHandler *handlers.ArticleHandler
}

func NewContainer() *Container {
	infraFactory := db.InfraFactory
	repos := infraFactory.CreateRepositoryPorts()

	// Initialize services
	authService := services.NewAuthService(repos.AdminRepo, repos.TokenRepo)
	userService := services.NewUserApplicationService(repos.UserRepo, authService)
	deviceService := services.NewDeviceApplicationService(repos.DeviceRepo, authService)
	articleService := services.NewArticleApplicationService(repos.ArticleRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)
	articleHandler := handlers.NewArticleHandler(articleService)

	return &Container{
		infraFactory:   infraFactory,
		authService:    authService,
		userService:    userService,
		deviceService:  deviceService,
		articleService: articleService,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		DeviceHandler:  deviceHandler,
		ArticleHandler: articleHandler,
	}
}
