package core

import (
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/rest"
	"github.com/ichaly/yugong/core/serv"
	"github.com/ichaly/yugong/core/serv/douyin"
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
		serv.NewQueue,
		douyin.NewScript,
		douyin.NewSpider,
		fx.Annotated{
			Group:  "plugin",
			Target: rest.NewDouyinApi,
		},
	),
	fx.Invoke(base.Bootstrap),
)
