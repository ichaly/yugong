package test

import (
	"github.com/ichaly/yugong/core/serv/douyin"
	"testing"
)

var dyScript *douyin.Script

func init() {
	dyScript, _ = douyin.NewScript()
}

func TestDyScript(t *testing.T) {
	t.Logf("dyScript init:%+v", dyScript)
}
