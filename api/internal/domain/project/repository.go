package project

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, project *Project) (*Project, error)
	List(ctx context.Context, search string, sort string, favoriteOnly bool) ([]*Project, int, error)
	GetById(ctx context.Context, id string) (*Project, error)
	Update(ctx context.Context, project *Project) (*Project, error)
	SetFavorite(ctx context.Context, id string, favorite bool) (*Project, error)
	Delete(ctx context.Context, id string) error
}
