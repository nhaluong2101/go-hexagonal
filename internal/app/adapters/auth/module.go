package auth

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"auth-handler-module",
	TokenModule,
)
