package project

import (
	"github.com/analopesdev/project-manager-clicksign/api/internal/domain/project"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotate(
		NewProjectRepository,
		fx.As(new(project.Repository)),
	),
)
