package matchers

import (
	"github.com/golang/glog"
	"regexp"
)

type Matcher interface {
	GetName() string
	IsMatch(src map[string]string) bool
}

type MatcherConfig struct {
	MatchField []MatchField `yaml:"matchField"`
}

type MatchField struct {
	Name       string  `yaml:"name"`
	FieldKey   string  `yaml:"fieldKey"`
	MatchStr   string  `yaml:"matchStr"`
	MatchRegex MyRegex `yaml:"matchRegex"`
}

func (t MatchField) GetName() string {
	return t.Name
}
func (t MatchField) IsMatch(src map[string]string) bool {

	val, ok := src[t.FieldKey]
	if !ok {
		return false
	}

	if t.MatchStr != "" {
		return t.MatchStr == val
	}

	return t.MatchRegex.IsMatch(val)
}

type MyRegex struct {
	RegexStr string         `yaml:"regex"`
	Regex    *regexp.Regexp `yaml:"-"`
}

func (t MyRegex) IsMatch(value string) bool {
	return t.Regex.MatchString(value)
}

func (t MyRegex) UnmarshalYAML(unmarshal func(interface{}) error) error {

	err := unmarshal(&t.RegexStr)
	if err != nil {
		return err
	}

	t.Regex, err = regexp.Compile(t.RegexStr)
	if err != nil {
		glog.Error("invalid regex %s, %s. ", t.RegexStr, err.Error())
		return err
	}
	return nil
}

func (t MyRegex) MarshalYAML() (interface{}, error) {
	return t.RegexStr, nil
}
