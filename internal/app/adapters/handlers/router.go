package handlers

import (
	"github.com/bagashiz/go_hexagonal/internal/app/core/ports"
	"github.com/bagashiz/go_hexagonal/internal/app/infrastructure/configs"
	"go.uber.org/fx"
	"log/slog"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/samber/slog-gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// RouterHandler is a wrapper for HTTP router
type RouterHandler struct {
	*gin.Engine
}

// NewRouterHandler creates a new HTTP router
func NewRouterHandler(
	config *configs.Container,
	token ports.TokenService,
	userHandler *UserHandler,
	authHandler *AuthHandler,
) (*RouterHandler, error) {

	// Disable debug mode in production
	if config.HTTP.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS
	ginConfig := cors.DefaultConfig()
	allowedOrigins := config.HTTP.AllowedOrigins
	originsList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originsList

	router := gin.New()
	router.Use(sloggin.New(slog.Default()), gin.Recovery(), cors.New(ginConfig))

	// Custom validators
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		if err := v.RegisterValidation("user_role", userRoleValidator); err != nil {
			return nil, err
		}

	}

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/", userHandler.Register)
			user.POST("/login", authHandler.Login)

			authUser := user.Group("/").Use(authMiddleware(token))
			{
				authUser.GET("/", userHandler.ListUsers)
				authUser.GET("/:id", userHandler.GetUser)

				admin := authUser.Use(adminMiddleware())
				{
					admin.PUT("/:id", userHandler.UpdateUser)
					admin.DELETE("/:id", userHandler.DeleteUser)
				}
			}
		}
	}

	return &RouterHandler{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *RouterHandler) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}

var RouterModule = fx.Module(
	"router-handler-module",
	fx.Provide(NewRouterHandler),
)
