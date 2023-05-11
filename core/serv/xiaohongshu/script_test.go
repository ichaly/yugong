package xiaohongshu

import (
	"testing"
)

func TestSign(t *testing.T) {
	s, err := NewScript()
	if err != nil {
		panic(err)
	}
	token := s.Sign("/api/sns/web/v1/user_posted?num=30&cursor=&user_id=55727cfbf5a2635fc1a9d345", nil)
	xt := token["X-t"]
	xs := token["X-s"]
	t.Logf("xt:%s >>> xs:%s", xt, xs)
	c := s.Common("18805433dffvrl8b2hut9pmkcdwizrt067cr7g6i330000326172", xt, xs)
	t.Log(c)
}
