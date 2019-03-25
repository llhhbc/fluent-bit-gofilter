package config

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func TestParseConfig(t *testing.T) {

	pc := ParseConfig{}

	src := `
match:
  matchField:
  - name: ma
`
	err := yaml.Unmarshal([]byte(src), &pc)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("unmarshal ok ", pc)
}
