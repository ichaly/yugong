package core

import (
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/plug"
	"github.com/ichaly/yugong/core/serv"
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
	fx.Provide(
		serv.NewSpider,
		fx.Annotated{
			Group:  "plugin",
			Target: plug.NewUserService,
		},
	),
	fx.Invoke(base.Bootstrap),
)
