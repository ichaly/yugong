package xiaohongshu

import (
	_ "embed"
	"github.com/dop251/goja"
)

//go:embed lib/origin_script.js
var bogus string

type Script struct {
	sign func(string, interface{}) map[string]string
	vm   *goja.Runtime
}

func NewScript() (*Script, error) {
	vm := goja.New()
	_, err := vm.RunString(bogus)
	if err != nil {
		return nil, err
	}
	var fn func(string, interface{}) map[string]string
	err = vm.ExportTo(vm.Get("sign"), &fn)
	if err != nil {
		return nil, err
	}
	return &Script{fn, vm}, nil
}

func (my *Script) Sign(query string, data map[string]string) map[string]string {
	if data == nil {
		return my.sign(query, nil)
	}
	return my.sign(query, my.vm.ToValue(data))
}
