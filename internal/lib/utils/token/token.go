package token

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("MY_4W3S0M3_S3CR37")

type contextKey string

const userIDContextKey contextKey = "user_id"

func Generate(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Срок действия — 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})
}

func UserIDFrom(ctx context.Context) (int, error) {
	user := ctx.Value(userIDContextKey)
	if user == nil {
		return 0, fmt.Errorf("cannot get user id from context")
	}

	userID, ok := user.(int)
	if !ok {
		return 0, fmt.Errorf("invalid user id type %T", user)
	}

	return userID, nil
}

func ContextWithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDFromClaims(claims jwt.MapClaims) (int, error) {
	rawUserID, ok := claims["user_id"]
	if !ok {
		return 0, fmt.Errorf("user_id claim is missing")
	}

	switch value := rawUserID.(type) {
	case float64:
		return int(value), nil
	case int:
		return value, nil
	case int64:
		return int(value), nil
	default:
		return 0, fmt.Errorf("invalid user_id claim type %T", rawUserID)
	}
}
