package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/analopesdev/project-manager-clicksign/api/internal/domain/project"
	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/httpapi/controllers"
	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/httpapi/routers"
	"github.com/go-chi/chi/v5"
)

const (
	knownID   = "3873c8e3-d436-4967-83b9-0a3daadd17eb"
	unknownID = "99999999-9999-4999-8999-999999999999"
	validBody = `{"name":"Redesign do app","client":"Clicksign","start_date":"2026-01-15","end_date":"2026-06-30"}`
)

// memoryRepository is a project.Repository kept in memory; err, when set, fails every call.
type memoryRepository struct {
	projects  map[string]*project.Project
	err       error
	listCalls []string
}

func (r *memoryRepository) Create(_ context.Context, p *project.Project) (*project.Project, error) {
	if r.err != nil {
		return nil, r.err
	}
	created := *p
	created.ID = unknownID[:len(unknownID)-1] + "0"
	r.projects[created.ID] = &created
	return &created, nil
}

func (r *memoryRepository) List(_ context.Context, search, sort string, favoriteOnly bool) ([]*project.Project, int, error) {
	if r.err != nil {
		return nil, 0, r.err
	}
	r.listCalls = append(r.listCalls, search+"|"+sort)
	list := []*project.Project{}
	for _, p := range r.projects {
		list = append(list, p)
	}
	return list, len(list), nil
}

func (r *memoryRepository) GetById(_ context.Context, id string) (*project.Project, error) {
	if r.err != nil {
		return nil, r.err
	}
	if p, ok := r.projects[id]; ok {
		copy := *p
		return &copy, nil
	}
	return nil, project.ErrProjectNotFound
}

func (r *memoryRepository) Update(_ context.Context, p *project.Project) (*project.Project, error) {
	r.projects[p.ID] = p
	return p, r.err
}

func (r *memoryRepository) SetFavorite(_ context.Context, id string, favorite bool) (*project.Project, error) {
	p, err := r.GetById(context.Background(), id)
	if err != nil {
		return nil, err
	}
	p.Favorite = favorite
	return p, nil
}

func (r *memoryRepository) Delete(_ context.Context, id string) error {
	if r.err != nil {
		return r.err
	}
	if _, ok := r.projects[id]; !ok {
		return project.ErrProjectNotFound
	}
	delete(r.projects, id)
	return nil
}

func newServer(t *testing.T) (*chi.Mux, *memoryRepository) {
	t.Helper()
	repo := &memoryRepository{projects: map[string]*project.Project{
		knownID: {ID: knownID, Name: "Original", Client: "Acme", StartDate: "2026-01-01", EndDate: "2026-02-01"},
	}}
	controller := controllers.NewProjectController(project.NewProjectService(repo), slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := chi.NewRouter()
	routers.ProjectRouter(router, controller)
	return router, repo
}

func call(router http.Handler, method, path, body string) (int, map[string]any) {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var out map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	return rec.Code, out
}

func TestProjectsHTTPContract(t *testing.T) {
	t.Run("creates a project", func(t *testing.T) {
		router, _ := newServer(t)
		status, body := call(router, http.MethodPost, "/projects/", validBody)
		if status != http.StatusCreated || body["name"] != "Redesign do app" {
			t.Fatalf("status = %d, body = %v", status, body)
		}
	})

	t.Run("an empty body is a 422 with one message per required field", func(t *testing.T) {
		router, _ := newServer(t)
		status, body := call(router, http.MethodPost, "/projects/", `{}`)
		fields, _ := body["fields"].(map[string]any)
		if status != http.StatusUnprocessableEntity || fields["name"] == nil || fields["client"] == nil || fields["start_date"] == nil || fields["end_date"] == nil {
			t.Fatalf("status = %d, body = %v", status, body)
		}
	})

	t.Run("an unknown field is a 400", func(t *testing.T) {
		router, _ := newServer(t)
		if status, _ := call(router, http.MethodPost, "/projects/", `{"nome":"campo desconhecido"}`); status != http.StatusBadRequest {
			t.Fatalf("status = %d", status)
		}
	})

	t.Run("lists as {data, total}", func(t *testing.T) {
		router, repo := newServer(t)
		status, body := call(router, http.MethodGet, "/projects/?search=acme&sort=end_date", "")
		data, _ := body["data"].([]any)
		if status != http.StatusOK || len(data) != 1 || body["total"] != float64(1) {
			t.Fatalf("status = %d, body = %v", status, body)
		}
		if len(repo.listCalls) != 1 || repo.listCalls[0] != "acme|end_date" {
			t.Fatalf("search and sort were not forwarded: %v", repo.listCalls)
		}
	})

	t.Run("an unknown sort is a 422 and never reaches the database", func(t *testing.T) {
		router, repo := newServer(t)
		for _, sort := range []string{"invalid", "name;DROP TABLE projects"} {
			status, body := call(router, http.MethodGet, "/projects/?sort="+url.QueryEscape(sort), "")
			fields, _ := body["fields"].(map[string]any)
			if status != http.StatusUnprocessableEntity || fields["sort"] == nil {
				t.Fatalf("sort %q: status = %d, body = %v", sort, status, body)
			}
		}
		if len(repo.listCalls) != 0 {
			t.Fatalf("the repository was queried: %v", repo.listCalls)
		}
	})

	t.Run("an unknown or malformed id is a 404 on every route", func(t *testing.T) {
		router, _ := newServer(t)
		for _, id := range []string{unknownID, "not-a-uuid"} {
			requests := []struct{ method, path, body string }{
				{http.MethodGet, "/projects/" + id, ""},
				{http.MethodPut, "/projects/" + id, validBody},
				{http.MethodPatch, "/projects/" + id + "/favorite", `{"favorite":true}`},
				{http.MethodDelete, "/projects/" + id, ""},
			}
			for _, rq := range requests {
				if status, _ := call(router, rq.method, rq.path, rq.body); status != http.StatusNotFound {
					t.Errorf("%s %s: status = %d, want 404", rq.method, rq.path, status)
				}
			}
		}
	})

	t.Run("deletes a project", func(t *testing.T) {
		router, _ := newServer(t)
		if status, _ := call(router, http.MethodDelete, "/projects/"+knownID, ""); status != http.StatusNoContent {
			t.Fatalf("delete: status = %d", status)
		}
		if status, _ := call(router, http.MethodGet, "/projects/"+knownID, ""); status != http.StatusNotFound {
			t.Fatalf("get after delete: status = %d", status)
		}
	})

	t.Run("a database failure is a generic 500 that leaks nothing", func(t *testing.T) {
		router, repo := newServer(t)
		repo.err = errors.New(`password authentication failed for user "postgres" (SQLSTATE 28P01)`)

		status, body := call(router, http.MethodGet, "/projects/", "")
		if status != http.StatusInternalServerError || body["error"] != "internal server error" {
			t.Fatalf("status = %d, body = %v", status, body)
		}
	})
}
