package config_test

import (
	"testing"

	"github.com/sanzhang007/crawl_nodes/config"
)

func TestLoadConfig(t *testing.T) {
	m, err := config.LoadConfig("../static/url.json")
	t.Logf("urls:%v,[err]:%v\n", m.Urls, err)
}
