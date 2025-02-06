package handlers

import (
	"github.com/bagashiz/go_hexagonal/internal/app/adapters/author"
	_constant "github.com/bagashiz/go_hexagonal/internal/app/core/constant"
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
	"github.com/bagashiz/go_hexagonal/internal/app/core/ports"
	"github.com/bagashiz/go_hexagonal/internal/app/core/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// TokenMiddleware is a author to check if the user is authenticated
func TokenMiddleware(token ports.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(_constant.AuthorizationHeaderKey)

		isEmpty := len(authorizationHeader) == 0
		if isEmpty {
			err := models.ErrEmptyAuthorizationHeader
			utils.HandleAbort(ctx, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		isValid := len(fields) == 2
		if !isValid {
			err := models.ErrInvalidAuthorizationHeader
			utils.HandleAbort(ctx, err)
			return
		}

		currentAuthorizationType := strings.ToLower(fields[0])
		if currentAuthorizationType != _constant.AuthorizationType {
			err := models.ErrInvalidAuthorizationType
			utils.HandleAbort(ctx, err)
			return
		}

		accessToken := fields[1]
		payload, err := token.VerifyToken(accessToken)
		if err != nil {
			utils.HandleAbort(ctx, err)
			return
		}

		ctx.Set(_constant.AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func RoleMiddleware(casbin *author.CasbinConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		payload := GetAuthPayload(ctx, _constant.AuthorizationPayloadKey)
		sub := payload.Role
		obj := ctx.Request.URL.Path
		act := ctx.Request.Method

		allowed, err := casbin.Enforcer.Enforce(sub, obj, act)
		if err != nil {
			err := models.ErrInternal
			utils.HandleAbort(ctx, err)
			return
		}
		if !allowed {
			err := models.ErrForbidden
			utils.HandleAbort(ctx, err)
			return
		}

		ctx.Next()
	}
}
