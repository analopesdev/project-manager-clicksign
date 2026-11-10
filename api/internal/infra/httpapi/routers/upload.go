package routers

import (
	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/httpapi/controllers"
	"github.com/go-chi/chi/v5"
)

func UploadRouter(r *chi.Mux, controller *controllers.UploadController) {
	r.Post("/uploads", controller.Upload)
	r.Get("/uploads/{name}", controller.Serve)
}
