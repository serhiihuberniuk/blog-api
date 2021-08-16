package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/serhiihuberniuk/blog-api/models"
	"golang.org/x/crypto/bcrypt"
)

type tokenClaims struct {
	jwt.StandardClaims
	userID string
}

func (s *Service) Login(ctx context.Context, payload models.LoginPayload) (string, error) {
	user, err := s.repo.Login(ctx, payload.Email)
	if err != nil {
		return "", models.ErrNotAuthenticated
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", models.ErrNotAuthenticated
		}

		return "", fmt.Errorf("error occurred while checking the password, %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID: user.ID,
	})

	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("error occurred while signing token: %w", err)
	}

	return tokenString, nil
}
