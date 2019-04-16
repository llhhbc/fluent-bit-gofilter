package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {

	pc := ParseConfig{}

	src := `
match:
  matchField:
  - name: ma
    matchStr: adfasd
parser:
  parseRegex:
  - name: pr1
    regex: sdafasdf
  parseJson:
  - name: pj1
    key: pj1
`
	err := yaml.UnmarshalStrict([]byte(src), &pc)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pc.MatcherConfig, pc.ParserConfig)
	err = pc.ApplyConfig()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("unmarshal ok %#v. ", pc)
}

func TestConfig(t *testing.T) {

	msg, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	c := Config{}
	err = yaml.UnmarshalStrict(msg, &c)
	if err != nil {
		t.Fatal(err)
	}

	err = c.ApplyConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("unmarshal ok %#v. ", c)
}

func TestTime(t *testing.T) {

	//tm, err := time.Parse("2006-01-02T15:04:05.999999+00:00", "2019-03-27T17:36:30.113317+08:00")
	tm, err := time.Parse(time.RFC3339Nano, "2019-03-27T17:36:30.113317+08:00")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("get tm:", tm, time.Now())
}
