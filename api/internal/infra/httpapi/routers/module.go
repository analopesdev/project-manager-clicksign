package routers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewRouter),
	fx.Invoke(ProjectRouter),
)
