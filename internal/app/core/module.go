package core

import (
	"github.com/bagashiz/go_hexagonal/internal/app/core/services"
	"go.uber.org/fx"
)

var Module = fx.Options(
	services.Module,
)
