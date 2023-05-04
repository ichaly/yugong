package main

import (
	"github.com/ichaly/yugong/core"
	"go.uber.org/fx"
	"path/filepath"
)

func main() {
	fx.New(
		core.Modules,
		fx.Provide(
			fx.Annotated{
				Name:   "configFile",
				Target: func() string { return filepath.Join("../conf", "dev.yml") },
			},
		),
	).Run()
}
