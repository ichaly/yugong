package data

import (
	"github.com/ichaly/yugong/core/util"
	"time"
)

type History struct {
	Title string     `jsonschema:"title=标题"`
	Type  string     `jsonschema:"title=类型"`
	Size  int64      `jsonschema:"title=大小"`
	Done  int64      `jsonschema:"title=完成大小"`
	Start time.Time  `jsonschema:"title=开始时间"`
	End   *time.Time `jsonschema:"title=完成时间"`
}

func (my *History) Progress() int64 {
	if my.Size == 0 {
		return 0
	}
	return my.Done * 100 / my.Size
}

func (my *History) Spend() time.Duration {
	if my.End == nil {
		my.End = util.TimePtr(time.Now())
	}
	return my.End.Sub(my.Start)
}
