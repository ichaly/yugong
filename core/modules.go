package core

import (
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/rest"
	"github.com/ichaly/yugong/core/serv"
	"github.com/ichaly/yugong/core/serv/douyin"
	"github.com/ichaly/yugong/core/serv/xiaohongshu"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(
		fx.Annotate(
			base.NewConfig,
			fx.ParamTags(`name:"configFile"`),
		),
		base.NewStore,
		base.NewCache,
		base.NewConnect,
		base.NewRender,
		base.NewRouter,
		base.NewServer,
	),
	fx.Provide(
		serv.NewQueue,
		serv.NewCrontab,
		douyin.NewScript,
		xiaohongshu.NewScript,
		fx.Annotated{
			Group:  "plugin",
			Target: rest.NewBlogApi,
		},
		fx.Annotate(
			douyin.NewDouyin,
			fx.As(new(serv.Spider)),
			fx.ResultTags(`group:"spider"`),
		),
		fx.Annotate(
			xiaohongshu.NewXiaoHongShu,
			fx.As(new(serv.Spider)),
			fx.ResultTags(`group:"spider"`),
		),
	),
	fx.Invoke(base.Bootstrap),
)
