package douyin

import (
	_ "embed"
	"github.com/dop251/goja"
)

//go:embed lib/X-Bogus.min.js
var bogus string

type Script struct {
	sign func(string, string) string
}

func NewScript() (*Script, error) {
	vm := goja.New()
	_, err := vm.RunString(bogus)
	if err != nil {
		return nil, err
	}
	var fn func(string, string) string
	err = vm.ExportTo(vm.Get("sign"), &fn)
	if err != nil {
		return nil, err
	}
	return &Script{fn}, nil
}

func (my *Script) Sign(query, agent string) string {
	return my.sign(query, agent)
}
