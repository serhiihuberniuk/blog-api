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
			return nil, fmt.Errorf("invalid signing method: %w", models.ErrNotAuthenticated)
		}

		return s.privateKey.Public(), nil
	})
	if err != nil {
		return "", fmt.Errorf("error ocured while parsing token: %w", models.ErrNotAuthenticated)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", fmt.Errorf("token claims are not of type *tokenClaims; %w", models.ErrNotAuthenticated)
	}

	return fmt.Sprint(claims.UserID), nil
}

func (s *Service) authPostAuthor(ctx context.Context, postID string) error {
	post, err := s.GetPost(ctx, postID)
	if err != nil {
		return fmt.Errorf("cannot get post: %w", err)
	}

	if post.CreatedBy != s.GetCurrentUserID(ctx) {
		return models.ErrNotAuthenticated
	}

	return nil
}

func (s *Service) authCommentAuthor(ctx context.Context, commentID string) error {
	post, err := s.GetComment(ctx, commentID)
	if err != nil {
		return fmt.Errorf("cannot get comment: %w", err)
	}

	if post.CreatedBy != s.GetCurrentUserID(ctx) {
		return models.ErrNotAuthenticated
	}

	return nil
}
