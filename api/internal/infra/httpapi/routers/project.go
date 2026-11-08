package routers

import (
	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/httpapi/controllers"
	"github.com/go-chi/chi/v5"
)

func ProjectRouter(r *chi.Mux, controller *controllers.ProjectController) {
	r.Route("/projects", func(r chi.Router) {
		r.Get("/", controller.ListProjects)
		r.Post("/", controller.CreateProject)
		r.Get("/{id}", controller.GetProjectById)
		r.Put("/{id}", controller.UpdateProject)
		r.Patch("/{id}/favorite", controller.SetFavorite)
		r.Delete("/{id}", controller.DeleteProject)
	})
}
