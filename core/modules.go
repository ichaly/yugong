package core

import (
	"github.com/ichaly/jingwei/core/base"
	"github.com/ichaly/jingwei/core/serv"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(
		base.NewConfig,
		base.NewStore,
		base.NewCache,
		base.NewConnect,
		base.NewRender,
		base.NewRouter,
		base.NewServer,
	),
	fx.Provide(fx.Annotated{
		Group:  "plugin",
		Target: serv.NewUserService,
	}),
	fx.Invoke(base.Bootstrap),
)
