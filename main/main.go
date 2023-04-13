package main

import (
	"github.com/ichaly/jingwei/core"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		core.Modules,
	).Run()
}
