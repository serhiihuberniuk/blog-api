package decorator

import (
	"context"

	"github.com/serhiihuberniuk/blog-api/models"
)

func (d *RepositoryCacheDecorator) Login(ctx context.Context, email string) (*models.User, error) {
	return d.repository.Login(ctx, email)
}
