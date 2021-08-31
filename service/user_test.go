package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/serhiihuberniuk/blog-api/models"
	mock_service "github.com/serhiihuberniuk/blog-api/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_GetUser(t *testing.T) {
	t.Parallel()

	type mockBehavior func(r *mock_service.Mockrepository, ctx context.Context, userID string)

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
			mockBehavior: func(r *mock_service.Mockrepository, ctx context.Context, userID string) {
				r.EXPECT().GetUser(context.Background(), gomock.Eq(user.ID)).
					Return(user, nil)
			},
			expectedUser: user,
			errMessage:   "",
		},
		{
			name:     "Error in repository layer",
			inUserId: user.ID,
			inCtx:    context.Background(),
			mockBehavior: func(r *mock_service.Mockrepository, ctx context.Context, userID string) {
				r.EXPECT().GetUser(context.Background(), gomock.Eq(user.ID)).
					Return(nil, errors.New("cannot get user"))
			},
			expectedUser: nil,
			errMessage:   "cannot get user",
		},
		{
			name:     "User is not found",
			inUserId: "Invalid userID",
			inCtx:    context.Background(),
			mockBehavior: func(r *mock_service.Mockrepository, ctx context.Context, userID string) {
				r.EXPECT().GetUser(context.Background(), gomock.Eq("Invalid userID")).
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
			mockRepo := mock_service.NewMockrepository(ctrl)
			service := Service{
				repo: mockRepo,
			}

			tc.mockBehavior(mockRepo, tc.inCtx, tc.inUserId)

			gottenUser, err := service.GetUser(tc.inCtx, tc.inUserId)
			if tc.errMessage == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, gottenUser)

				return
			}

			assert.Contains(t, err.Error(), tc.errMessage)
			assert.Nil(t, gottenUser)
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	t.Parallel()

	type repoMockBehavior func(r *mock_service.Mockrepository, ctx context.Context, userId string)

	type providerMockBehavior func(p *mock_service.MockcurrentUserInformationProvider, ctx context.Context) string

	userID := "40d2f213-8e41-4abf-ad54-410997b19401"
	ctxWithUserId := context.WithValue(context.Background(), "userID", userID)

	testCases := []struct {
		name                 string
		inCtx                context.Context
		repoMockBehavior     repoMockBehavior
		providerMockBehavior providerMockBehavior
		errMessage           string
	}{
		{
			name:  "User is deleted",
			inCtx: ctxWithUserId,
			repoMockBehavior: func(r *mock_service.Mockrepository, ctx context.Context, userId string) {
				r.EXPECT().DeleteUser(gomock.Eq(ctxWithUserId), gomock.Eq(userID)).Return(nil)
			},
			providerMockBehavior: func(p *mock_service.MockcurrentUserInformationProvider, ctx context.Context) string {
				return p.EXPECT().GetCurrentUserID(gomock.Eq(ctxWithUserId)).Return(userID).String()
			},
			errMessage: "",
		},
		{
			name:  "Error in repository layer",
			inCtx: ctxWithUserId,
			repoMockBehavior: func(r *mock_service.Mockrepository, ctx context.Context, userId string) {
				r.EXPECT().DeleteUser(gomock.Eq(ctxWithUserId), gomock.Eq(userID)).
					Return(errors.New("cannot delete user"))
			},
			providerMockBehavior: func(p *mock_service.MockcurrentUserInformationProvider, ctx context.Context) string {
				return p.EXPECT().GetCurrentUserID(gomock.Eq(ctxWithUserId)).Return(userID).String()
			},
			errMessage: "cannot delete user",
		},
		{
			name:  "User is not found",
			inCtx: ctxWithUserId,
			repoMockBehavior: func(r *mock_service.Mockrepository, ctx context.Context, userId string) {
				r.EXPECT().DeleteUser(gomock.Eq(ctxWithUserId), gomock.Eq("Invalid userID")).Return(models.ErrNotFound)
			},
			providerMockBehavior: func(p *mock_service.MockcurrentUserInformationProvider, ctx context.Context) string {
				return p.EXPECT().GetCurrentUserID(gomock.Eq(ctxWithUserId)).Return("Invalid userID").String()
			},
			errMessage: models.ErrNotFound.Error(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repoMock := mock_service.NewMockrepository(ctrl)
			providerMock := mock_service.NewMockcurrentUserInformationProvider(ctrl)
			service := Service{
				repo:                           repoMock,
				currentUserInformationProvider: providerMock,
			}

			tc.repoMockBehavior(repoMock, tc.inCtx, tc.providerMockBehavior(providerMock, tc.inCtx))

			err := service.DeleteUser(tc.inCtx)
			if tc.errMessage == "" {
				assert.NoError(t, err)

				return
			}

			assert.Contains(t, err.Error(), tc.errMessage)
		})
	}
}
