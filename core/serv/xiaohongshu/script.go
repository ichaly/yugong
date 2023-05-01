package xiaohongshu

import (
	_ "embed"
	"github.com/dop251/goja"
)

//go:embed lib/origin_script.js
var bogus string

type Script struct {
	sign func(string, string) map[string]string
}

func NewScript() (*Script, error) {
	vm := goja.New()
	_, err := vm.RunString(bogus)
	if err != nil {
		return nil, err
	}
	var fn func(string, string) map[string]string
	err = vm.ExportTo(vm.Get("sign"), &fn)
	if err != nil {
		return nil, err
	}
	return &Script{fn}, nil
}

func (my *Script) Sign(query, data string) map[string]string {
	return my.sign(query, data)
}
