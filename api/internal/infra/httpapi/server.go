package httpapi

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/config"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

func NewServer(r *chi.Mux, cfg *config.Config) *http.Server {
	return &http.Server{Addr: ":" + cfg.Port, Handler: r}
}

func RunServer(lc fx.Lifecycle, srv *http.Server, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			logger.Info("Server is running", "addr", srv.Addr)
			go func() {
				if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("Server stopped", "err", err)
				}
			}()
			return nil
		},
		OnStop: srv.Shutdown,
	})
}
