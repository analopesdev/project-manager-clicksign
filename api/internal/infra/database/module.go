package database

import (
	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/database/project"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewPostgresDB,
	),
	project.Module,
)
