package service_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/service"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Login(t *testing.T) {
	t.Parallel()

	type mockBehavior func(r *Mockrepository, ctx context.Context, email string, user *models.User)
	input := &models.LoginPayload{
		Email:    "email@mail.com",
		Password: "password",
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Log(fmt.Errorf("error occurred while running test: %w", err))
		t.Fail()

		return
	}

	user := &models.User{
		ID:        "315b6c09-36ff-4519-8579-492f3ae2a3be",
		Name:      "name",
		Email:     "email@email.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Password:  string(hashedPassword),
	}
	testCases := []struct {
		name         string
		inCtx        context.Context
		inPayload    models.LoginPayload
		mockBehavior mockBehavior
		expectedUser *models.User
		errMessage   string
	}{
		{
			name:  "Login OK",
			inCtx: context.Background(),
			inPayload: models.LoginPayload{
				Email:    input.Email,
				Password: input.Password,
			},
			mockBehavior: func(r *Mockrepository, ctx context.Context, email string, user *models.User) {
				r.EXPECT().Login(gomock.Eq(ctx), gomock.Eq(email)).Return(user, nil)
			},
			expectedUser: user,
			errMessage:   "",
		},
		{
			name:  "Repository error",
			inCtx: context.Background(),
			inPayload: models.LoginPayload{
				Email:    input.Email,
				Password: input.Password,
			},
			mockBehavior: func(r *Mockrepository, ctx context.Context, email string, user *models.User) {
				r.EXPECT().Login(gomock.Eq(ctx), gomock.Eq(email)).Return(nil, models.ErrNotAuthenticated)
			},
			errMessage: models.ErrNotAuthenticated.Error(),
		},
		{
			name:  "Invalid payload(no email)",
			inCtx: context.Background(),
			inPayload: models.LoginPayload{
				Password: input.Password,
			},
			mockBehavior: func(r *Mockrepository, ctx context.Context, email string, user *models.User) {
				r.EXPECT().Login(gomock.Eq(ctx), gomock.Eq(email)).Return(nil, models.ErrNotAuthenticated)
			},
			errMessage: models.ErrNotAuthenticated.Error(),
		},
		{
			name:  "Invalid password",
			inCtx: context.Background(),
			inPayload: models.LoginPayload{
				Email:    input.Email,
				Password: "Invalid password",
			},
			mockBehavior: func(r *Mockrepository, ctx context.Context, email string, user *models.User) {
				r.EXPECT().Login(gomock.Eq(ctx), gomock.Eq(email)).Return(user, nil)
			},
			expectedUser: user,
			errMessage:   models.ErrNotAuthenticated.Error(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := NewMockrepository(ctrl)
			providerMock := NewMockcurrentUserInformationProvider(ctrl)

			privatKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				t.Log(fmt.Errorf("error occurred while generating privat key: %w", err))
				t.Fail()

				return
			}

			serv, err := service.NewService(repoMock, privatKey, providerMock)
			if err != nil {
				t.Log(fmt.Errorf("error occured while init service: %w", err))
				t.Fail()

				return
			}

			tc.mockBehavior(repoMock, tc.inCtx, tc.inPayload.Email, tc.expectedUser)

			token, err := serv.Login(tc.inCtx, tc.inPayload)
			if tc.errMessage == "" {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				return
			}

			assert.Empty(t, token)
			assert.Contains(t, err.Error(), tc.errMessage)

		})
	}
}
