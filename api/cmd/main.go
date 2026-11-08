package main

import (
	"log/slog"
	"os"

	"github.com/analopesdev/project-manager-clicksign/api/internal/domain/project"
	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/config"
	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/database"
	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/httpapi"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
)

func main() {

	app := fx.New(
		fx.Provide(
			func() *slog.Logger { return slog.New(slog.NewJSONHandler(os.Stdout, nil)) },
			config.Load,
		),
		database.Module,
		project.Module,
		httpapi.Module,
		fx.NopLogger,
	)

	app.Run()
}
