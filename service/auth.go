package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/serhiihuberniuk/blog-api/models"
	"time"
)

const key = "mysecretkey"

type tokenClaims struct {
	jwt.StandardClaims
	userID string
}

func (s *Service) Login(ctx context.Context, payload models.LoginPayload) (string, error) {

	user, err := s.repo.Login(ctx, payload.Email, payload.Password)
	fmt.Println(err)
	if err != nil {
		return "", fmt.Errorf("authentication failed :%w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID: user.ID,
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("authentification failed: %w", err)
	}

	return tokenString, nil
}
