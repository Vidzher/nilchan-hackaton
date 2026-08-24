package token

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("MY_4W3S0M3_S3CR37")

type contextKey string

const userIDContextKey contextKey = "user_id"

func Generate(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Срок действия — 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func signingKey(_ *jwt.Token) (any, error) {
	return secretKey, nil
}

func Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(
		tokenString,
		signingKey,
		jwt.WithValidMethods([]string{"HS256"}),
	)
}

func UserIDFromContext(ctx context.Context) (int64, error) {
	user := ctx.Value(userIDContextKey)
	if user == nil {
		return 0, fmt.Errorf("cannot get user id from context")
	}

	userID, ok := user.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid user id type %T", user)
	}

	return userID, nil
}

func ContextWithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDFromClaims(claims jwt.MapClaims) (int64, error) {
	subject, err := claims.GetSubject() //subject - встроенный в jwt "владелец" токена
	if err != nil {
		return 0, fmt.Errorf("get subject claim: %w", err)
	}
	if subject == "" {
		return 0, fmt.Errorf("subject claim is missing")
	}

	userID, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid subject claim: %w", err)
	}

	return userID, nil
}
