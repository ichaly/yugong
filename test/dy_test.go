package test

import (
	"github.com/ichaly/yugong/core/serv/douyin"
	"github.com/ichaly/yugong/core/util"
	"testing"
	"time"
)

var dy *douyin.Douyin

func init() {
	dy = douyin.NewDouyin(connect, config, dyScript)
}

func TestDouyin(t *testing.T) {
	err := dy.GetVideos(
		"MS4wLjABAAAAuHyViCTGymyGLtUl_G0y_rh2jVaObXbF8t7a7GeFDRw",
		"2209433555684",
		nil,
		util.StringPtr(util.FormatLong(time.Now().UnixMilli())),
		nil,
		5,
		0,
	)
	if err != nil {
		panic(err)
	}
}
