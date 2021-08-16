package service

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
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
		return "", fmt.Errorf("error occurred while getting user :%w", err)
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

	privateKey, err := ioutil.ReadFile(s.config.PrivateKeyPath)
	if err != nil {
		return "", fmt.Errorf("error occurred while reading private.pem file :%w", err)
	}

	privateRSA, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", fmt.Errorf("error occurred while parsing private key: %w", err)
	}

	tokenString, err := token.SignedString(privateRSA)
	if err != nil {
		return "", fmt.Errorf("error occurred while signing token: %w", err)
	}

	return tokenString, nil
}
