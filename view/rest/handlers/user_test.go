package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang/mock/gomock"
	"github.com/serhiihuberniuk/blog-api/models"
	"github.com/serhiihuberniuk/blog-api/view/rest/handlers"
	"github.com/stretchr/testify/assert"
)

func TestHandlers_CreateUser(t *testing.T) {
	t.Parallel()

	type serviceCreateUserMockBehavior func(s *Mockservice, ctx context.Context, payload models.CreateUserPayload,
		userID string)
	type serviceGetUserMockBehavior func(s *Mockservice, ctx context.Context, userID string, user *models.User)

	input := models.CreateUserPayload{
		Name:     "name",
		Email:    "email",
		Password: "password",
	}

	inputBody := `{"name": "name", "email": "email", "password": "password"}`
	expectedBody := `{"id":"id","name":"name","email":"email","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z"}
`
	testCases := []struct {
		name                          string
		inputBody                     string
		inputPayload                  models.CreateUserPayload
		serviceCreateUserMockBehavior serviceCreateUserMockBehavior
		gotUserIdFromCreateUser       string
		serviceGetUserMockBehavior    serviceGetUserMockBehavior
		gotUser                       *models.User
		expectedCode                  int
		expectedBody                  string
	}{
		{
			name:         "Status OK",
			inputBody:    inputBody,
			inputPayload: input,

			serviceCreateUserMockBehavior: func(s *Mockservice, ctx context.Context, payload models.CreateUserPayload, userID string) {
				s.EXPECT().CreateUser(ctx, payload).Return(userID, nil)
			},
			gotUserIdFromCreateUser: "id",
			serviceGetUserMockBehavior: func(s *Mockservice, ctx context.Context, userID string, user *models.User) {
				s.EXPECT().GetUser(ctx, userID).Return(user, nil)
			},
			gotUser: &models.User{
				ID:    "id",
				Name:  input.Name,
				Email: input.Email,
			},
			expectedCode: http.StatusOK,
			expectedBody: expectedBody,
		},
		{
			name:      "Bad request",
			inputBody: `{"email":"email", "password":"password"}`,
			inputPayload: models.CreateUserPayload{
				Email:    "email",
				Password: "password",
			},
			serviceCreateUserMockBehavior: func(s *Mockservice, ctx context.Context, payload models.CreateUserPayload, userID string) {
				s.EXPECT().CreateUser(ctx, gomock.Any()).Return("", validation.Errors{})
			},
			gotUserIdFromCreateUser:    "",
			serviceGetUserMockBehavior: func(s *Mockservice, ctx context.Context, userID string, user *models.User) {},
			gotUser:                    &models.User{},
			expectedCode:               http.StatusBadRequest,
			expectedBody:               "Bad Request\n",
		},
		{
			name:         "Internal server error",
			inputBody:    inputBody,
			inputPayload: input,
			serviceCreateUserMockBehavior: func(s *Mockservice, ctx context.Context, payload models.CreateUserPayload, userID string) {
				s.EXPECT().CreateUser(ctx, payload).Return(userID, nil)
			},
			gotUserIdFromCreateUser: "id",
			serviceGetUserMockBehavior: func(s *Mockservice, ctx context.Context, userID string, user *models.User) {
				s.EXPECT().GetUser(ctx, userID).Return(nil, errors.New("internal error"))
			},
			gotUser:      nil,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Internal Server Error\n",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			servMock := NewMockservice(ctrl)
			providerMock := NewMockcurrentUserInformationProvider(ctrl)
			middlewareMock := NewMockauthMiddleware(ctrl)
			handlersRest := handlers.NewRestHandlers(servMock, middlewareMock, providerMock)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/users", bytes.NewBufferString(tc.inputBody))
			tc.serviceCreateUserMockBehavior(servMock, r.Context(), tc.inputPayload, tc.gotUserIdFromCreateUser)
			tc.serviceGetUserMockBehavior(servMock, r.Context(), tc.gotUserIdFromCreateUser, tc.gotUser)

			handler := http.HandlerFunc(handlersRest.CreateUser)
			handler.ServeHTTP(w, r)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
