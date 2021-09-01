package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/service"
	"github.com/stretchr/testify/assert"
)

func TestService_GetUser(t *testing.T) {
	t.Parallel()

	type mockBehavior func(r *Mockrepository, ctx context.Context, userID string, user *models.User)

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
			mockBehavior: func(r *Mockrepository, ctx context.Context, userID string, user *models.User) {
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
			mockBehavior: func(r *Mockrepository, ctx context.Context, userID string, user *models.User) {
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
			mockBehavior: func(r *Mockrepository, ctx context.Context, userID string, user *models.User) {
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
			defer ctrl.Finish()
			mockRepo := NewMockrepository(ctrl)
			mockProvider := NewMockcurrentUserInformationProvider(ctrl)

			serv, err := service.NewService(mockRepo, nil, mockProvider)
			if err != nil {
				t.Log(err)
				t.Fail()

				return
			}

			tc.mockBehavior(mockRepo, tc.inCtx, tc.inUserId, tc.expectedUser)

			gotUser, err := serv.GetUser(tc.inCtx, tc.inUserId)
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

	type repoMockBehavior func(r *Mockrepository, ctx context.Context, userId string)

	type providerMockBehavior func(p *MockcurrentUserInformationProvider, ctx context.Context, userId string)

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
			repoMockBehavior: func(r *Mockrepository, ctx context.Context, userId string) {
				r.EXPECT().DeleteUser(gomock.Eq(ctx), gomock.Eq(userId)).Return(nil)
			},
			providerMockBehavior: func(p *MockcurrentUserInformationProvider, ctx context.Context, userId string) {
				p.EXPECT().GetCurrentUserID(gomock.Eq(ctxWithUserId)).Return(userId).AnyTimes()
			},
			expectedUserIdFromProvider: validUserId,
			errMessage:                 "",
		},
		{
			name:  "Error in repository layer",
			inCtx: ctxWithUserId,
			repoMockBehavior: func(r *Mockrepository, ctx context.Context, userId string) {
				r.EXPECT().DeleteUser(gomock.Eq(ctx), gomock.Eq(userId)).
					Return(errors.New("cannot delete user"))
			},
			providerMockBehavior: func(p *MockcurrentUserInformationProvider, ctx context.Context, userId string) {
				p.EXPECT().GetCurrentUserID(gomock.Eq(ctx)).Return(userId).AnyTimes()
			},
			expectedUserIdFromProvider: validUserId,
			errMessage:                 "cannot delete user",
		},
		{
			name:  "User is not found",
			inCtx: ctxWithUserId,
			repoMockBehavior: func(r *Mockrepository, ctx context.Context, userId string) {
				r.EXPECT().DeleteUser(gomock.Eq(ctx), gomock.Eq(userId)).Return(models.ErrNotFound)
			},
			providerMockBehavior: func(p *MockcurrentUserInformationProvider, ctx context.Context, userId string) {
				p.EXPECT().GetCurrentUserID(gomock.Eq(ctx)).Return(userId).AnyTimes()
			},
			expectedUserIdFromProvider: "Invalid userID",
			errMessage:                 models.ErrNotFound.Error(),
		},
		{
			name:  "Context is empty",
			inCtx: nil,
			repoMockBehavior: func(r *Mockrepository, ctx context.Context, userId string) {
				r.EXPECT().DeleteUser(gomock.Eq(ctx), gomock.Eq(userId)).Return(models.ErrNotFound)
			},
			providerMockBehavior: func(p *MockcurrentUserInformationProvider, ctx context.Context, userId string) {
				p.EXPECT().GetCurrentUserID(gomock.Eq(ctx)).Return(userId).Times(2)
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
			defer ctrl.Finish()
			repoMock := NewMockrepository(ctrl)
			providerMock := NewMockcurrentUserInformationProvider(ctrl)
			serv, err := service.NewService(repoMock, nil, providerMock)
			if err != nil {
				t.Log(fmt.Errorf("error occured while initialization of service: %w", err))
				t.Fail()

				return
			}
			tc.providerMockBehavior(providerMock, tc.inCtx, tc.expectedUserIdFromProvider)
			tc.repoMockBehavior(repoMock, tc.inCtx, providerMock.GetCurrentUserID(tc.inCtx))

			err = serv.DeleteUser(tc.inCtx)
			if tc.errMessage == "" {
				assert.NoError(t, err)

				return
			}

			assert.Contains(t, err.Error(), tc.errMessage)
		})
	}
}

func TestService_CreateUser(t *testing.T) {
	t.Parallel()

	input := &models.User{
		Name:     "name",
		Email:    "email@mail.com",
		Password: "password",
	}

	type mockBehavior func(r *Mockrepository, ctx context.Context)

	testCases := []struct {
		name         string
		inCtx        context.Context
		inPayload    models.CreateUserPayload
		mockBehavior mockBehavior
		expectedUser *models.User
		errMessage   string
	}{
		{
			name:  "User created",
			inCtx: context.Background(),
			inPayload: models.CreateUserPayload{
				Name:     input.Name,
				Email:    input.Email,
				Password: input.Password,
			},
			mockBehavior: func(r *Mockrepository, ctx context.Context) {
				r.EXPECT().CreateUser(gomock.Eq(ctx), gomock.Any()).
					Return(nil)
			},
			errMessage: "",
		},
		{
			name:  "Error in repository layer",
			inCtx: context.Background(),
			inPayload: models.CreateUserPayload{
				Name:     input.Name,
				Email:    input.Email,
				Password: input.Password,
			},
			mockBehavior: func(r *Mockrepository, ctx context.Context) {
				r.EXPECT().CreateUser(gomock.Eq(ctx), gomock.Any()).
					Return(errors.New("cannot create user"))
			},
			errMessage: "cannot create user",
		},
		{
			name:  "Invalid payload(name is empty)",
			inCtx: context.Background(),
			inPayload: models.CreateUserPayload{
				Email:    input.Email,
				Password: input.Password,
			},
			mockBehavior: func(r *Mockrepository, ctx context.Context) {},
			errMessage:   "validation failed",
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

			serv, err := service.NewService(repoMock, nil, providerMock)
			if err != nil {
				t.Log(fmt.Errorf("error occured while initialization of service: %w", err))
				t.Fail()

				return
			}

			tc.mockBehavior(repoMock, tc.inCtx)
			userID, err := serv.CreateUser(tc.inCtx, tc.inPayload)
			if tc.errMessage == "" {
				assert.NoError(t, err)
				assert.NotEmpty(t, userID)

				return
			}
			assert.Contains(t, err.Error(), tc.errMessage)
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	t.Parallel()

	validUserId := "315b6c09-36ff-4519-8579-492f3ae2a3be"
	input := &models.User{
		Name:     "name",
		Email:    "email@email.com",
		Password: "password",
	}
	ctxWithUserId := context.WithValue(context.Background(), "userID", validUserId)

	type repoMockBehaviorUpdate func(r *Mockrepository, ctx context.Context, userId string, user *models.User)
	type providerMockBehavior func(r *MockcurrentUserInformationProvider, ctx context.Context, userID string)

	testCases := []struct {
		name                       string
		inCtx                      context.Context
		inPayload                  models.UpdateUserPayload
		repoMockBehaviorUpdate     repoMockBehaviorUpdate
		providerMockBehavior       providerMockBehavior
		expectedUserIdFromProvider string
		expectedUserFromGetUser    *models.User
		errMessage                 string
	}{
		{
			name:  "User updated",
			inCtx: ctxWithUserId,
			inPayload: models.UpdateUserPayload{
				Name:     input.Name,
				Email:    input.Email,
				Password: input.Password,
			},
			providerMockBehavior: func(r *MockcurrentUserInformationProvider, ctx context.Context, userID string) {
				r.EXPECT().GetCurrentUserID(gomock.Eq(ctx)).Return(userID).Times(2)
			},
			repoMockBehaviorUpdate: func(r *Mockrepository, ctx context.Context, userID string, user *models.User) {
				r.EXPECT().GetUser(gomock.Eq(ctx), gomock.Eq(userID)).Return(user, nil)
				r.EXPECT().UpdateUser(gomock.Eq(ctx), gomock.Any()).Return(nil)
			},
			expectedUserFromGetUser:    input,
			expectedUserIdFromProvider: validUserId,
			errMessage:                 "",
		},
		{
			name:  "Context is empty",
			inCtx: nil,
			inPayload: models.UpdateUserPayload{
				Name:     input.Name,
				Email:    input.Email,
				Password: input.Password,
			},
			providerMockBehavior: func(r *MockcurrentUserInformationProvider, ctx context.Context, userID string) {
				r.EXPECT().GetCurrentUserID(gomock.Eq(ctx)).Return("").Times(2)
			},
			repoMockBehaviorUpdate: func(r *Mockrepository, ctx context.Context, userID string, user *models.User) {
				r.EXPECT().GetUser(gomock.Eq(ctx), gomock.Eq(userID)).Return(nil, models.ErrNotFound)
			},
			expectedUserFromGetUser:    nil,
			expectedUserIdFromProvider: "",
			errMessage:                 models.ErrNotFound.Error(),
		},
		{
			name:  "Repository error while getting user",
			inCtx: ctxWithUserId,
			inPayload: models.UpdateUserPayload{
				Name:     input.Name,
				Email:    input.Email,
				Password: input.Password,
			},
			providerMockBehavior: func(r *MockcurrentUserInformationProvider, ctx context.Context, userID string) {
				r.EXPECT().GetCurrentUserID(gomock.Eq(ctx)).Return(userID).Times(2)
			},
			repoMockBehaviorUpdate: func(r *Mockrepository, ctx context.Context, userID string, user *models.User) {
				r.EXPECT().GetUser(gomock.Eq(ctx), gomock.Eq(userID)).
					Return(nil, errors.New("cannot get user"))
			},
			expectedUserFromGetUser:    nil,
			expectedUserIdFromProvider: input.ID,
			errMessage:                 "cannot get user",
		},
		{
			name:  "Invalid payload(email empty)",
			inCtx: ctxWithUserId,
			inPayload: models.UpdateUserPayload{
				Name:     input.Name,
				Password: input.Password,
			},
			providerMockBehavior: func(r *MockcurrentUserInformationProvider, ctx context.Context, userID string) {
				r.EXPECT().GetCurrentUserID(gomock.Eq(ctx)).Return(userID).Times(2)
			},
			repoMockBehaviorUpdate: func(r *Mockrepository, ctx context.Context, userID string, user *models.User) {
				r.EXPECT().GetUser(gomock.Eq(ctx), gomock.Eq(userID)).Return(user, nil)
			},
			expectedUserFromGetUser:    input,
			expectedUserIdFromProvider: input.ID,
			errMessage:                 "validation failed",
		},
		{
			name:  "Repository error while updating user",
			inCtx: ctxWithUserId,
			inPayload: models.UpdateUserPayload{
				Name:     input.Name,
				Email:    input.Email,
				Password: input.Password,
			},
			providerMockBehavior: func(r *MockcurrentUserInformationProvider, ctx context.Context, userID string) {
				r.EXPECT().GetCurrentUserID(gomock.Eq(ctx)).Return(userID).Times(2)
			},
			repoMockBehaviorUpdate: func(r *Mockrepository, ctx context.Context, userID string, user *models.User) {
				r.EXPECT().GetUser(gomock.Eq(ctx), gomock.Eq(userID)).Return(user, nil)
				r.EXPECT().UpdateUser(gomock.Eq(ctx), gomock.Any()).Return(errors.New("cannot update user"))
			},
			expectedUserFromGetUser:    input,
			expectedUserIdFromProvider: input.ID,
			errMessage:                 "cannot update user",
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

			serv, err := service.NewService(repoMock, nil, providerMock)
			if err != nil {
				t.Log(fmt.Errorf("error occured while initialization of service: %w", err))
				t.Fail()

				return
			}

			tc.providerMockBehavior(providerMock, tc.inCtx, tc.expectedUserIdFromProvider)
			tc.repoMockBehaviorUpdate(repoMock, tc.inCtx, providerMock.GetCurrentUserID(tc.inCtx), tc.expectedUserFromGetUser)

			err = serv.UpdateUser(tc.inCtx, tc.inPayload)
			if tc.errMessage == "" {
				assert.NoError(t, err)

				return
			}

			assert.Contains(t, err.Error(), tc.errMessage)
		})
	}
}
