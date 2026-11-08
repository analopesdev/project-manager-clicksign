package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/analopesdev/project-manager-clicksign/api/internal/domain/project"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ProjectController struct {
	service *project.ProjectService
	logger  *slog.Logger
}

func NewProjectController(service *project.ProjectService, logger *slog.Logger) *ProjectController {
	return &ProjectController{service: service, logger: logger}
}

var validateProject = validator.New()

type CreateOrUpdateProjectRequest struct {
	Name      string `json:"name" validate:"required,max=120"`
	Client    string `json:"client" validate:"required,max=120"`
	StartDate string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"required,datetime=2006-01-02"`
	CoverURL  string `json:"cover_url" validate:"omitempty,url"`
}

func (r *CreateOrUpdateProjectRequest) Validate() error {
	return validateProject.Struct(r)
}

type ListProjectsResponse struct {
	Data  []*project.Project `json:"data"`
	Total int                `json:"total"`
}

type ListProjectsParams struct {
	Search       string `json:"search"`
	Sort         string `json:"sort"`
	FavoriteOnly bool   `json:"favorite"`
}

func (r *ListProjectsParams) ValidateListProjectsParams() error {
	return validateProject.Struct(r)
}

type SetFavoriteRequest struct {
	Favorite *bool `json:"favorite"`
}

func (c *ProjectController) ListProjects(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := ListProjectsParams{
		Search:       query.Get("search"),
		Sort:         query.Get("sort"),
		FavoriteOnly: query.Get("favorite") == "true" || query.Get("favorite") == "1",
	}

	result, total, err := c.service.ListProjects(r.Context(), params.Search, params.Sort, params.FavoriteOnly)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, ListProjectsResponse{Data: result, Total: total})
}

func (c *ProjectController) GetProjectById(w http.ResponseWriter, r *http.Request) {
	p, err := c.service.GetProjectById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: err.Error()})
		return
	}
	if err := writeJSON(w, http.StatusOK, p); err != nil {
		c.logger.ErrorContext(r.Context(), "failed to write response", "error", err)
		return
	}
}

func (c *ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	var dto CreateOrUpdateProjectRequest

	fmt.Println("dto", dto)

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := dto.Validate(); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, errorResponse{Error: err.Error()})
		return
	}

	p, err := c.service.CreateProject(r.Context(), dto.Name, dto.Client, dto.StartDate, dto.EndDate, dto.CoverURL)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (c *ProjectController) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var request CreateOrUpdateProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := request.Validate(); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, errorResponse{Error: err.Error()})
		return
	}

	p, err := c.service.UpdateProject(r.Context(), chi.URLParam(r, "id"), request.Name, request.Client, request.StartDate, request.EndDate, request.CoverURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (c *ProjectController) SetFavorite(w http.ResponseWriter, r *http.Request) {
	var dto SetFavoriteRequest

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if dto.Favorite == nil {
		writeJSON(w, http.StatusUnprocessableEntity, errorResponse{Error: "favorite is required"})
		return
	}

	p, err := c.service.SetFavorite(r.Context(), chi.URLParam(r, "id"), *dto.Favorite)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (c *ProjectController) DeleteProject(w http.ResponseWriter, r *http.Request) {
	if err := c.service.DeleteProject(r.Context(), chi.URLParam(r, "id")); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
