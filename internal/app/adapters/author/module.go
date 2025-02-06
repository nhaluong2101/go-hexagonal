package author

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"auth-handler-module",
	CasbinModule,
)
