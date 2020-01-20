package config

import (
	"github.com/beaquant/utils/json_file"
	"testing"
)

func TestConfig(t *testing.T) {
	c := &Config{}
	t.Log(json_file.Load("config-sample.json", c))
	t.Log(c)
}
