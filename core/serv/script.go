package serv

import (
	_ "embed"
	"github.com/dop251/goja"
)

var (
	//go:embed lib/X-Bogus.min.js
	bogus string
	agent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.39"
)

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

func (my *Script) Sign(query string) string {
	return my.sign(query, agent)
}

func (my *Script) Agent() string {
	return agent
}
