package handlers

import (
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/author"
	"github.com/bagashiz/go_hexagonal/internal/app/core/ports"
	"github.com/bagashiz/go_hexagonal/internal/app/infrastructure/configs"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/samber/slog-gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log/slog"
	"os"
	"strings"
	"time"
)

// RouterHandler is a wrapper for HTTP router
type RouterHandler struct {
	*gin.Engine
}

const logPath = "./logs/go.log"

var logger *zap.Logger

func setupLog() {
	_, err := os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout", logPath}
	logger, _ = c.Build()
}

// NewRouterHandler creates a new HTTP router
func NewRouterHandler(
	config *configs.Container,
	casbin *author.CasbinConfig,
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

	//POLICY
	casbin.LoadPolicy()

	//LOGGER
	setupLog()
	// Setting GIN to use zap as logger
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	v1 := router.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/", userHandler.Register)
			user.POST("/login", authHandler.Login)

			authUser := user.Group("/").Use(TokenMiddleware(token), RoleMiddleware(casbin))
			{
				authUser.GET("/", userHandler.ListUsers)
				authUser.GET("/:id", userHandler.GetUser)

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
