package xiaohongshu

import (
	_ "embed"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/dop251/goja"
	"github.com/ichaly/yugong/core/util"
	"strings"
)

var (
	index int
	//go:embed lib/sign.js
	sign string
	//go:embed lib/mcr.js
	mcr string
	//go:embed lib/commonId.js
	common string
)

type Script struct {
	svm    *goja.Runtime
	mvm    *goja.Runtime
	cvm    *goja.Runtime
	sign   func(string, interface{}) map[string]string
	common func(encode []byte) string
	mcr    func(string) int64
}

func NewScript() (*Script, error) {
	s := Script{
		svm: goja.New(),
		mvm: goja.New(),
		cvm: goja.New(),
	}

	_, err := s.svm.RunString(sign)
	if err != nil {
		return nil, err
	}
	err = s.svm.ExportTo(s.svm.Get("sign"), &s.sign)
	if err != nil {
		return nil, err
	}

	_, err = s.mvm.RunString(mcr)
	if err != nil {
		return nil, err
	}
	err = s.mvm.ExportTo(s.mvm.Get("mcrFun"), &s.mcr)
	if err != nil {
		return nil, err
	}

	_, err = s.cvm.RunString(common)
	if err != nil {
		return nil, err
	}
	err = s.cvm.ExportTo(s.cvm.Get("b64Encode"), &s.common)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (my *Script) Sign(query string, data map[string]string) map[string]string {
	if data == nil {
		return my.sign(query, nil)
	}
	return my.sign(query, my.svm.ToValue(data))
}

func (my *Script) Common(str, xt, xs string) string {
	x8 := "I38rHdgsjopgIvesdVwgIC+oIELmBZ5e3VwXLgFTIxS3bqwErFeexd0ekncAzMFYnqthIhJeSnMDKutRI3KsYorWHPtGrbV0P9WfIi/eWc6eYqtyQApPI37ekmR1QL+5Ii6sdnosjoT5yqtXqqwYrBqoIx++GDi/sVtkIx0sxuwr4qtiIkrwIi/skcc3ICLfI3Oe0utl20DZsL5eDSJejVw0IieexVwL+PtorqthPleekmW4Ix+iIhrRzVwgKPtYI3QPorKe6qthIx/s6VtoIiMoIiEge96eSdckrfvsjutKrZgefcr9gfKeYcPsIxKefVtzIE6edutholFIQdLnIx5s3qtRnc7eVfG+IkcwIiHt40bMIvhqtS8gIiifpVwAICHVJo3eSpeekPtVIx3e0jH="
	args := map[string]any{
		"s0":  5,
		"s1":  "",
		"x0":  "1",
		"x1":  "3.2.0",
		"x2":  "Windows",
		"x3":  "xhs-pc-web",
		"x4":  "2.0.3",
		"x5":  str, //cookie_a1
		"x6":  util.ParseLong(xt),
		"x7":  xs,
		"x8":  x8,
		"x9":  my.mcr(fmt.Sprintf("%s%s%s", xt, xs, x8)),
		"x10": index,
	}
	json, _ := sonic.MarshalString(args)
	json = strings.ReplaceAll(json, " ", "")
	index += 1
	return my.common(encodeUtf8(json))
}

func encodeUtf8(str string) []byte {
	utf8Bytes := []byte(str)
	encodedBytes := make([]byte, 0, len(utf8Bytes))
	for _, b := range utf8Bytes {
		if b < 0x80 {
			encodedBytes = append(encodedBytes, b)
		} else {
			encodedBytes = append(encodedBytes, 0xc0|(b>>6))
			encodedBytes = append(encodedBytes, 0x80|(b&0x3f))
		}
	}
	return encodedBytes
}
