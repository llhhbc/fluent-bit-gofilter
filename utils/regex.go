package utils

import (
	"github.com/golang/glog"
	"regexp"
)

type MyRegex struct {
	RegexStr string         `yaml:"regex"`
	Regex    *regexp.Regexp `yaml:"-"`
}

func (t MyRegex) IsMatch(value string) bool {
	return t.Regex.MatchString(value)
}

func (t *MyRegex) UnmarshalYAML(unmarshal func(interface{}) error) error {

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
