package main

import (
	"github.com/ichaly/yugong/core"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		core.Modules,
	).Run()
}
