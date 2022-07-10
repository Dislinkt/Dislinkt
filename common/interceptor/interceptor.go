package interceptor

import (
	"context"
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	accessiblePermissions map[string]string
	publicKey             string
}

func NewAuthInterceptor(accessiblePermissions map[string]string, publicKey string) *AuthInterceptor {
	return &AuthInterceptor{
		accessiblePermissions: accessiblePermissions,
		publicKey:             publicKey,
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fmt.Println(info.FullMethod)
		ctx, err := interceptor.Authorize(ctx, info.FullMethod)
		if err != nil {
			fmt.Println("EROOOOOOOOOOOOOOOOOOOOOOR")
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) Authorize(ctx context.Context, method string) (context.Context, error) {
	fmt.Println("USAOOOOOOOOOO")
	accessiblePermission, ok := interceptor.accessiblePermissions[method]
	// u mapi ne postoje role za ovu metodu => javno dostupna putanja
	if !ok {
		return ctx, nil
	}

	var values []string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("NEMA METADATA")
		return ctx, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	values = md.Get("Authorization")
	if len(values) == 0 {
		fmt.Println("NEMA AUTH")
		return ctx, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	authHeader := values[0]
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		fmt.Println("NIJE SPLIT")
		return ctx, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	claims, err := interceptor.verifyToken(parts[1])
	if err != nil {
		fmt.Println("NIJE VALIDAN")
		return ctx, status.Errorf(codes.Unauthenticated, "Unauthorized")
	}

	for _, role := range claims.Permissions {
		if role == accessiblePermission {
			fmt.Println(role)
			return context.WithValue(ctx, LoggedInUserKey{}, claims.Username), nil
		}
	}

	return ctx, status.Errorf(codes.PermissionDenied, "Forbidden")
}

func (interceptor *AuthInterceptor) verifyToken(accessToken string) (claims *Claims, err error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(interceptor.publicKey), nil
	})

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, fmt.Errorf("Couldn't parse claims")
	}

	if !claims.VerifyExpiresAt(time.Now().Local().Unix(), true) {
		return nil, fmt.Errorf("JWT is expired")
	}

	return claims, nil
}

type Claims struct {
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

type LoggedInUserKey struct{}
