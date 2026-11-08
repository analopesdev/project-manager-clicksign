package httpapi

import (
	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/httpapi/controllers"
	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/httpapi/routers"
	"go.uber.org/fx"
)

var Module = fx.Options(
	controllers.Module,
	routers.Module,
	fx.Provide(NewServer),
	fx.Invoke(RunServer),
)
