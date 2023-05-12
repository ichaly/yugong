package test

import (
	"github.com/ichaly/yugong/core/base"
	"path/filepath"
	"testing"
)

var config *base.Config

func init() {
	config, _ = base.NewConfig(filepath.Join("../conf", "dev.yml"))
}

func TestConfig(t *testing.T) {
	t.Logf("config init:%+v", config)
}
