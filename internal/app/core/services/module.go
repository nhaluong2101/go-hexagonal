package services

import (
	"github.com/bagashiz/go_hexagonal/internal/app/core/ports"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	fx.Provide(
		fx.Annotate(NewUserService, fx.As(new(ports.UserService))),
		fx.Annotate(NewAuthService, fx.As(new(ports.AuthService))),
	),
)
