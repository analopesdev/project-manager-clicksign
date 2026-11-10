package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/analopesdev/project-manager-clicksign/api/internal/domain/project"
	"github.com/go-chi/chi/v5"
)

const maxBodyBytes = 1 << 20

type ProjectController struct {
	service *project.ProjectService
	logger  *slog.Logger
}

func NewProjectController(service *project.ProjectService, logger *slog.Logger) *ProjectController {
	return &ProjectController{service: service, logger: logger}
}

type CreateOrUpdateProjectRequest struct {
	Name      string `json:"name"`
	Client    string `json:"client"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	CoverURL  string `json:"cover_url"`
}

type ListProjectsResponse struct {
	Data  []*project.Project `json:"data"`
	Total int                `json:"total"`
}

type ListProjectsParams struct {
	Search       string
	Sort         string
	FavoriteOnly bool
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
		c.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, ListProjectsResponse{Data: result, Total: total})
}

func (c *ProjectController) GetProjectById(w http.ResponseWriter, r *http.Request) {
	p, err := c.service.GetProjectById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (c *ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	var dto CreateOrUpdateProjectRequest
	if !decodeJSON(w, r, &dto) {
		return
	}

	p, err := c.service.CreateProject(r.Context(), dto.Name, dto.Client, dto.StartDate, dto.EndDate, dto.CoverURL)
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (c *ProjectController) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var dto CreateOrUpdateProjectRequest
	if !decodeJSON(w, r, &dto) {
		return
	}

	p, err := c.service.UpdateProject(r.Context(), chi.URLParam(r, "id"), dto.Name, dto.Client, dto.StartDate, dto.EndDate, dto.CoverURL)
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (c *ProjectController) SetFavorite(w http.ResponseWriter, r *http.Request) {
	var dto SetFavoriteRequest
	if !decodeJSON(w, r, &dto) {
		return
	}
	if dto.Favorite == nil {
		writeJSON(w, http.StatusUnprocessableEntity, validationResponse{
			Error:  "validation failed",
			Fields: map[string]string{"favorite": "favorite is required"},
		})
		return
	}

	p, err := c.service.SetFavorite(r.Context(), chi.URLParam(r, "id"), *dto.Favorite)
	if err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (c *ProjectController) DeleteProject(w http.ResponseWriter, r *http.Request) {
	if err := c.service.DeleteProject(r.Context(), chi.URLParam(r, "id")); err != nil {
		c.writeServiceError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// decodeJSON reads the body into dst, refusing unknown fields (400) and bodies over 1 MiB (413).
// It answers the request itself and returns false when the body is not usable.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, errors.New("request body must be at most 1 MiB"))
			return false
		}
		writeError(w, http.StatusBadRequest, err)
		return false
	}
	return true
}

// writeServiceError turns a domain error into its HTTP answer. Anything unexpected is logged and
// answered with a generic 500, so database details never reach the client.
func (c *ProjectController) writeServiceError(w http.ResponseWriter, r *http.Request, err error) {
	var validation *project.ValidationError
	switch {
	case errors.As(err, &validation):
		writeJSON(w, http.StatusUnprocessableEntity, validationResponse{Error: validation.Error(), Fields: validation.Fields})
	case errors.Is(err, project.ErrProjectNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: project.ErrProjectNotFound.Error()})
	default:
		c.logger.ErrorContext(r.Context(), "request failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}
