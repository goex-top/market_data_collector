package config

import (
	"github.com/jinzhu/configor"
	"testing"
)

func TestConfig(t *testing.T) {
	c1 := &Config{}
	t.Log(configor.Load(c1, "config-sample.json"))
	t.Log(c1)
	c2 := &Config{}
	t.Log(configor.Load(c2, "config-sample.yml"))
	t.Log(c2)
}
