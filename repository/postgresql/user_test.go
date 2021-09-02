package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/serhiihuberniuk/blog-api/models"
	repository "github.com/serhiihuberniuk/blog-api/repository/postgresql"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func createPostgresTestContainer(ctx context.Context) (func(), *repository.Repository, error) {
	dbName := "test_postgres_db"
	port := "5432/tcp"
	env := map[string]string{
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_DB":       dbName,
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres",
			Env:          env,
			ExposedPorts: []string{port},
			WaitingFor:   wait.ForListeningPort(nat.Port(port)),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("error occurred while creating container: %w", err)
	}

	terminateContainer := func() {
		err = container.Terminate(ctx)
		if err != nil {
			panic(err)
		}
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return terminateContainer, nil, fmt.Errorf("error occurred while getting port: %w", err)
	}

	dbUrl := fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", mappedPort.Port(), dbName)

	pool, err := repository.NewPostgresDb(ctx, dbUrl, "../../init.sql")
	if err != nil {
		return terminateContainer, nil, fmt.Errorf("cannot create conn pool: %w", err)
	}

	repo := &repository.Repository{
		Db: pool,
	}

	return terminateContainer, repo, nil
}

func cleanUserTable(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
	t.Helper()

	_, err := pool.Exec(ctx, "delete from users")
	if err != nil {
		t.Log(fmt.Errorf("cannot clean users table: %w", err))
	}
}

func TestRepository_User(t *testing.T) {
	ctx := context.Background()

	terminateContainer, repo, err := createPostgresTestContainer(ctx)
	if err != nil {
		t.Log(fmt.Errorf("error occurred while creating Postgres test container: %w", err))
		t.Fail()

		return
	}

	defer terminateContainer()
	defer repo.Db.Close()

	creationTime := time.Now()
	updatingTime := time.Now().Add(time.Minute)
	userToCreate := &models.User{
		ID:        "a209a589-4559-42eb-aab2-0863fbaa5f45",
		Name:      "name",
		Email:     "email@mail.com",
		CreatedAt: creationTime,
		UpdatedAt: creationTime,
		Password:  "hashedPassword",
	}

	userToUpdate := &models.User{
		ID:        userToCreate.ID,
		Name:      "UpdatedName",
		Email:     "UpdatedEmail@mail.com",
		CreatedAt: creationTime,
		UpdatedAt: updatingTime,
		Password:  "updatedHashedPassword",
	}
	notValidUserID := "b209a589-4559-42eb-aab2-0863fbaa5f45"

	testCases := []struct {
		name                string
		userToCreate        *models.User
		userIdToGetUser     string
		userIdToDelete      string
		userToUpdate        *models.User
		expectedUser        *models.User
		expectedUpdatedUser *models.User
		errMessage          string
	}{
		{
			name:                "User created, updated and deleted",
			userToCreate:        userToCreate,
			userIdToGetUser:     userToCreate.ID,
			userIdToDelete:      userToCreate.ID,
			userToUpdate:        userToUpdate,
			expectedUser:        userToCreate,
			expectedUpdatedUser: userToUpdate,
			errMessage:          "",
		},
		{
			name: "UserID is not provided",
			userToCreate: &models.User{
				Name:      userToCreate.Name,
				Email:     userToCreate.Email,
				CreatedAt: userToCreate.CreatedAt,
				UpdatedAt: userToCreate.UpdatedAt,
				Password:  userToCreate.Password,
			},
			errMessage: "invalid input syntax for type uuid",
		},
		{
			name:            "Cannot find user with such id",
			userToCreate:    userToCreate,
			userIdToGetUser: notValidUserID,
			errMessage:      models.ErrNotFound.Error(),
		},
		{
			name:            "Cannot find user to update",
			userToCreate:    userToCreate,
			userIdToGetUser: userToCreate.ID,
			userToUpdate: &models.User{
				ID:        notValidUserID,
				Name:      userToCreate.Name,
				Email:     userToCreate.Email,
				CreatedAt: userToCreate.CreatedAt,
				UpdatedAt: userToCreate.UpdatedAt,
				Password:  userToCreate.Password,
			},
			expectedUser: userToCreate,
			errMessage:   models.ErrNotFound.Error(),
		},
		{
			name:                "Cannot delete user",
			userToCreate:        userToCreate,
			userIdToGetUser:     userToCreate.ID,
			userIdToDelete:      notValidUserID,
			userToUpdate:        userToUpdate,
			expectedUser:        userToCreate,
			expectedUpdatedUser: userToUpdate,
			errMessage:          models.ErrNotFound.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer cleanUserTable(t, repo.Db, ctx)

			err = repo.CreateUser(ctx, tc.userToCreate)
			if err != nil {
				if tc.errMessage != "" {
					assert.Contains(t, err.Error(), tc.errMessage)

					return
				}
				t.Fail()

				return
			}
			assert.NoError(t, err)

			gotUser, err := repo.GetUser(ctx, tc.userIdToGetUser)
			if err != nil {
				if tc.errMessage != "" {
					assert.Contains(t, err.Error(), tc.errMessage)
					assert.Empty(t, gotUser)

					return
				}
				t.Fail()

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUser.ID, gotUser.ID)
			assert.Equal(t, tc.expectedUser.Name, gotUser.Name)
			assert.Equal(t, tc.expectedUser.Email, gotUser.Email)

			err = repo.UpdateUser(ctx, tc.userToUpdate)
			if err != nil {
				if tc.errMessage != "" {
					assert.Contains(t, err.Error(), tc.errMessage)

					return
				}
				t.Fail()

				return
			}
			assert.NoError(t, err)

			updatedUser, err := repo.GetUser(ctx, userToUpdate.ID)
			if err != nil {
				if tc.errMessage != "" {
					assert.Contains(t, err.Error(), tc.errMessage)
					assert.Empty(t, gotUser)

					return
				}
				t.Fail()

				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUpdatedUser.ID, updatedUser.ID)
			assert.Equal(t, tc.expectedUpdatedUser.Name, updatedUser.Name)
			assert.Equal(t, tc.expectedUpdatedUser.Email, updatedUser.Email)
			assert.NotEqual(t, gotUser, updatedUser)

			err = repo.DeleteUser(ctx, tc.userIdToDelete)
			if err != nil {
				if tc.errMessage != "" {
					assert.Contains(t, err.Error(), tc.errMessage)

					return
				}
			}
			assert.NoError(t, err)

			user, err := repo.GetUser(ctx, tc.userIdToDelete)

			assert.Error(t, err)
			assert.Empty(t, user)
		})
	}
}
