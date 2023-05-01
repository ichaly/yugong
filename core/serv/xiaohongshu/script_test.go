package xiaohongshu

import (
	"testing"
)

func TestSign(t *testing.T) {
	s, err := NewScript()
	if err != nil {
		panic(err)
	}
	v := s.Sign("/api/sns/web/v1/user_posted?num=30&cursor=&user_id=55727cfbf5a2635fc1a9d345")
	t.Log(v)
}
