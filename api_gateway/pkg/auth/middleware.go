package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	pb "github.com/dislinkt/common/proto/auth_service"
	"github.com/gin-gonic/gin"
)

type AuthMiddlewareConfig struct {
	svc *ServiceClient
}

func InitAuthMiddleware(svc *ServiceClient) AuthMiddlewareConfig {
	return AuthMiddlewareConfig{svc}
}

func (config *AuthMiddlewareConfig) AuthRequired(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("authorization")

	if authorization == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	res, err := config.svc.Client.ValidateToken(context.Background(), &pb.ValidateRequest{
		Jwt: &pb.JwtToken{Jwt: token[1]},
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("role", res.Role)
	ctx.Next()
}

func (config *AuthMiddlewareConfig) Authorize(ctx *gin.Context) {
	val, existed := ctx.Get("role")
	fmt.Println(val.(string))
	accessible := AccessibleRoles()
	accessibleRoles, ok := accessible[ctx.FullPath()]
	// u mapi ne postoje role za ovu metodu => javno dostupna putanja
	if !ok {
		ctx.Next()
	}

	if !existed {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	for _, role := range accessibleRoles {
		if role == "ok" {
			ctx.Next()
		}
	}

	ctx.AbortWithStatus(http.StatusForbidden)
	return

}
