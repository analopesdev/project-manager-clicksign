package project

import (
	"context"
	"regexp"
)

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// A malformed id cannot exist, so it is "not found" without asking the database.
func validID(id string) bool { return uuidPattern.MatchString(id) }

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
	if err := ValidateSort(sort); err != nil {
		return nil, 0, err
	}
	return s.repository.List(ctx, search, sort, favoriteOnly)
}

func (s *ProjectService) GetProjectById(ctx context.Context, id string) (*Project, error) {
	if !validID(id) {
		return nil, ErrProjectNotFound
	}
	return s.repository.GetById(ctx, id)
}

func (s *ProjectService) UpdateProject(ctx context.Context, id string, name string, client string, startDate string, endDate string, coverURL string) (*Project, error) {
	if err := ValidateProjectInput(name, client, startDate, endDate, coverURL); err != nil {
		return nil, err
	}

	if !validID(id) {
		return nil, ErrProjectNotFound
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
	if !validID(id) {
		return nil, ErrProjectNotFound
	}
	return s.repository.SetFavorite(ctx, id, favorite)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	if !validID(id) {
		return ErrProjectNotFound
	}
	return s.repository.Delete(ctx, id)
}
