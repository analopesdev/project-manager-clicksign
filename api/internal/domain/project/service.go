package project

import (
	"context"
)

type ProjectService struct {
	repository Repository
}

func NewProjectService(repository Repository) *ProjectService {
	return &ProjectService{repository: repository}
}

func (s *ProjectService) CreateProject(ctx context.Context, name string, client string, startDate string, endDate string, coverURL string) (*Project, error) {
	if err := ValidateProjectInput(name, client, startDate, endDate, coverURL); err != nil {
		return nil, err
	}

	p := &Project{
		Name:      name,
		Client:    client,
		StartDate: startDate,
		EndDate:   endDate,
		CoverURL:  coverURL,
	}

	return s.repository.Create(ctx, p)
}

func (s *ProjectService) ListProjects(ctx context.Context, search string, sort string, favoriteOnly bool) ([]*Project, int, error) {
	return s.repository.List(ctx, search, sort, favoriteOnly)
}

func (s *ProjectService) GetProjectById(ctx context.Context, id string) (*Project, error) {
	return s.repository.GetById(ctx, id)
}

func (s *ProjectService) UpdateProject(ctx context.Context, id string, name string, client string, startDate string, endDate string, coverURL string) (*Project, error) {
	if err := ValidateProjectInput(name, client, startDate, endDate, coverURL); err != nil {
		return nil, err
	}

	p, err := s.repository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	p.Name = name
	p.Client = client
	p.StartDate = startDate
	p.EndDate = endDate
	p.CoverURL = coverURL

	return s.repository.Update(ctx, p)
}

func (s *ProjectService) SetFavorite(ctx context.Context, id string, favorite bool) (*Project, error) {
	return s.repository.SetFavorite(ctx, id, favorite)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
