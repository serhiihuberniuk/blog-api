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
	UserID string `json:"userId"`
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

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.ID,
	})

	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("error occurred while signing token: %w", err)
	}

	return tokenString, nil
}

func (s *Service) ParseToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}

		return s.privateKey.Public(), nil
	})
	if err != nil {
		return "", fmt.Errorf("error occurred while parsing token: %w", err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", fmt.Errorf("token claims are not of type *tokenClaims; %w", models.ErrNotAuthenticated)
	}

	return claims.UserID, nil
}

func (s *Service) checkCurrentUserIsOwner(ctx context.Context, id string) bool {
	if id == s.currentUserInformationProvider.GetCurrentUserID(ctx) {
		return true
	}

	return false
}
