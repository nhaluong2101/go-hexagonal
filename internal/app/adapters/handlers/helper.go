package handlers

import (
	"github.com/bagashiz/go_hexagonal/internal/app/core/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// stringToUint64 is a helper function to convert a string to uint64
func stringToUint64(str string) (uint64, error) {
	num, err := strconv.ParseUint(str, 10, 64)

	return num, err
}

// getAuthPayload is a helper function to get the auth payload from the context
func getAuthPayload(ctx *gin.Context, key string) *models.TokenPayload {
	return ctx.MustGet(key).(*models.TokenPayload)
}

// toMap is a helper function to add meta and data to a map
func toMap(m meta, data any, key string) map[string]any {
	return map[string]any{
		"meta": m,
		key:    data,
	}
}
