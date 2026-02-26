package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractTokenFromHeader(ctx *gin.Context, header string) string {
	authHeader := ctx.GetHeader("Authorization")
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}

	return ctx.Query("token")
}
