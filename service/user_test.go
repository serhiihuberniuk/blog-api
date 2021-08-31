package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/stretchr/testify/assert"
)

func TestService_GetUser(t *testing.T) {
	t.Parallel()

	type mockBehavior func(r *MockRepository, ctx context.Context, userID string, user *models.User)

	user := &models.User{
		ID:        "315b6c09-36ff-4519-8579-492f3ae2a3be",
		Name:      "name",
		Email:     "email@email.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Password:  "hashedPassword",
	}

	testCases := []struct {
		name         string
		inUserId     string
		inCtx        context.Context
		mockBehavior mockBehavior
		expectedUser *models.User
		errMessage   string
	}{
		{
			name:     "User is gotten",
			inUserId: user.ID,
			inCtx:    context.Background(),
			mockBehavior: func(r *MockRepository, ctx context.Context, userID string, user *models.User) {
				r.EXPECT().GetUser(gomock.Eq(ctx), gomock.Eq(userID)).
					Return(user, nil)
			},
			expectedUser: user,
			errMessage:   "",
		},
		{
			name:     "Error in repository layer",
			inUserId: user.ID,
			inCtx:    context.Background(),
			mockBehavior: func(r *MockRepository, ctx context.Context, userID string, user *models.User) {
				r.EXPECT().GetUser(gomock.Eq(ctx), gomock.Eq(userID)).
					Return(nil, errors.New("cannot get user"))
			},
			expectedUser: nil,
			errMessage:   "cannot get user",
		},
		{
			name:     "User is not found",
			inUserId: "Invalid userID",
			inCtx:    context.Background(),
			mockBehavior: func(r *MockRepository, ctx context.Context, userID string, user *models.User) {
				r.EXPECT().GetUser(gomock.Eq(ctx), gomock.Eq(userID)).
					Return(nil, models.ErrNotFound)
			},
			expectedUser: nil,
			errMessage:   models.ErrNotFound.Error(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockRepo := NewMockRepository(ctrl)
			service := Service{
				repo: mockRepo,
			}

			tc.mockBehavior(mockRepo, tc.inCtx, tc.inUserId, tc.expectedUser)

			gotUser, err := service.GetUser(tc.inCtx, tc.inUserId)
			if tc.errMessage == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, gotUser)

				return
			}

			assert.Contains(t, err.Error(), tc.errMessage)
			assert.Nil(t, gotUser)
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	t.Parallel()

	validUserId := "40d2f213-8e41-4abf-ad54-410997b19401"
	ctxWithUserId := context.WithValue(context.Background(), "userID", validUserId)

	type repoMockBehavior func(r *MockRepository, ctx context.Context, userId string)

	type providerMockBehavior func(p *MockCurrentUserInformationProvider, ctx context.Context, userId string)

	testCases := []struct {
		name                       string
		inCtx                      context.Context
		repoMockBehavior           repoMockBehavior
		providerMockBehavior       providerMockBehavior
		expectedUserIdFromProvider string
		errMessage                 string
	}{
		{
			name:  "User is deleted",
			inCtx: ctxWithUserId,
			repoMockBehavior: func(r *MockRepository, ctx context.Context, userId string) {
				r.EXPECT().DeleteUser(gomock.Eq(ctx), gomock.Eq(userId)).Return(nil)
			},
			providerMockBehavior: func(p *MockCurrentUserInformationProvider, ctx context.Context, userId string) {
				p.EXPECT().GetCurrentUserID(gomock.Eq(ctxWithUserId)).Return(userId).AnyTimes()
			},
			expectedUserIdFromProvider: validUserId,
			errMessage:                 "",
		},
		{
			name:  "Error in repository layer",
			inCtx: ctxWithUserId,
			repoMockBehavior: func(r *MockRepository, ctx context.Context, userId string) {
				r.EXPECT().DeleteUser(gomock.Eq(ctx), gomock.Eq(userId)).
					Return(errors.New("cannot delete user"))
			},
			providerMockBehavior: func(p *MockCurrentUserInformationProvider, ctx context.Context, userId string) {
				p.EXPECT().GetCurrentUserID(gomock.Eq(ctx)).Return(userId).AnyTimes()
			},
			expectedUserIdFromProvider: validUserId,
			errMessage:                 "cannot delete user",
		},
		{
			name:  "User is not found",
			inCtx: ctxWithUserId,
			repoMockBehavior: func(r *MockRepository, ctx context.Context, userId string) {
				r.EXPECT().DeleteUser(gomock.Eq(ctx), gomock.Eq(userId)).Return(models.ErrNotFound)
			},
			providerMockBehavior: func(p *MockCurrentUserInformationProvider, ctx context.Context, userId string) {
				p.EXPECT().GetCurrentUserID(gomock.Eq(ctx)).Return(userId).AnyTimes()
			},
			expectedUserIdFromProvider: "Invalid userID",
			errMessage:                 models.ErrNotFound.Error(),
		},
		{
			name:  "Context is empty",
			inCtx: nil,
			repoMockBehavior: func(r *MockRepository, ctx context.Context, userId string) {
				r.EXPECT().DeleteUser(gomock.Eq(ctx), gomock.Eq(userId)).Return(models.ErrNotFound)
			},
			providerMockBehavior: func(p *MockCurrentUserInformationProvider, ctx context.Context, userId string) {
				p.EXPECT().GetCurrentUserID(gomock.Eq(ctx)).Return(userId).AnyTimes()
			},
			expectedUserIdFromProvider: "",
			errMessage:                 models.ErrNotFound.Error(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repoMock := NewMockRepository(ctrl)
			providerMock := NewMockCurrentUserInformationProvider(ctrl)
			service := Service{
				repo:                           repoMock,
				currentUserInformationProvider: providerMock,
			}
			tc.providerMockBehavior(providerMock, tc.inCtx, tc.expectedUserIdFromProvider)
			tc.repoMockBehavior(repoMock, tc.inCtx, providerMock.GetCurrentUserID(tc.inCtx))

			err := service.DeleteUser(tc.inCtx)
			if tc.errMessage == "" {
				assert.NoError(t, err)

				return
			}

			assert.Contains(t, err.Error(), tc.errMessage)
		})
	}
}
