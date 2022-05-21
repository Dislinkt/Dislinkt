package interceptor

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	accessibleRoles map[string][]string
	publicKey       *rsa.PublicKey
}

func NewAuthInterceptor(accessibleRoles map[string][]string, publicKey *rsa.PublicKey) *AuthInterceptor {
	return &AuthInterceptor{
		accessibleRoles: accessibleRoles,
		publicKey:       publicKey,
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
	accessibleRoles, ok := interceptor.accessibleRoles[method]
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

	fmt.Println(claims["role"].(string))
	for _, role := range accessibleRoles {
		if role == claims["role"].(string) {
			fmt.Println(role)
			return context.WithValue(ctx, LoggedInUserKey{}, claims["username"].(string)), nil
		}
	}

	return ctx, status.Errorf(codes.PermissionDenied, "Forbidden")
}

func (interceptor *AuthInterceptor) verifyToken(accessToken string) (claims jwt.MapClaims, err error) {
	//token, err := jwt.ParseWithClaims(
	//	accessToken,
	//	&UserClaims{},
	//	func(token *jwt.Token) (interface{}, error) {
	//		_, ok := token.Method.(*jwt.SigningMethodRSA)
	//		if !ok {
	//			return nil, fmt.Errorf("Unexpected token signing method")
	//		}
	//
	//		return interceptor.publicKey, nil
	//	},
	//)
	//if err != nil {
	//	return nil, fmt.Errorf("Invalid token: %w", err)
	//}
	//claims, ok := token.Claims.(*UserClaims)
	//if !ok {
	//	return nil, fmt.Errorf("Invalid token claims")
	//}
	//
	//return claims, nil
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		//return []byte(os.Getenv("ACCESS_SECRET")), nil
		return []byte("Dislinkt"), nil
	})

	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("Couldn't parse claims")
	}

	//if !claims.VerifyExpiresAt(time.Now().Local().Unix()) {
	//	return nil, errors.New("JWT is expired")
	//}

	return claims, nil
}

type LoggedInUserKey struct{}
