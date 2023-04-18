package serv

import (
	"github.com/ichaly/yugong/core/util"
	"strings"
	"testing"
)

func TestSignature(t *testing.T) {
	in, err := util.ReadExcel("/Users/Chaly/Documents/Workspace/yugong/docs/header.xlsx", "params", "cookies")
	if err != nil {
		panic(err)
	}
	keys := map[string]bool{
		"appKey": true, "t": true, "sign": true, "data": true,
		//"cookie2": true,
		"_m_h5_tk": true, "_m_h5_tk_enc": true, "sg": true,
	}
	params := map[string]string{}
	cookies := map[string]string{}
	for _, p := range in["params"] {
		if _, ok := keys[p[0]]; ok {
			params[strings.Trim(p[0], "")] = strings.Trim(p[1], "")
			continue
		}
	}
	for _, c := range in["cookies"] {
		if _, ok := keys[c[0]]; ok {
			cookies[strings.Trim(c[0], "")] = strings.Trim(c[1], "")
			continue
		}
	}
	u := NewUploader()
	err = u.signature(params, cookies, 0)
	if err != nil {
		return
	}
}

func TestLogin(t *testing.T) {
	u := NewUploader()
	err := u.Login()
	if err != nil {
		return
	}
}
