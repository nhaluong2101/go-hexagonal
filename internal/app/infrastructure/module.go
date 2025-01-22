package infrastructure

import (
	"github.com/bagashiz/go_hexagonal/internal/app/infrastructure/configs"
	"go.uber.org/fx"
)

var Module = fx.Options(
	configs.Module,
)
