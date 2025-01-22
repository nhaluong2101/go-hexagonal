package handlers

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"handlers-module",
	UserModule,
	AuthModule,
	RouterModule,
)
