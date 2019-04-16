package matchers

import (
	"cgolib/utils"
	"github.com/golang/glog"
)

type Matcher interface {
	GetName() string
	IsMatch(src map[string]string) bool
}

type MatcherConfig struct {
	MatchField []*MatchField `yaml:"matchField"`
}

var _ Matcher = &MatchField{}

type MatchField struct {
	utils.Field `yaml:",inline"`
	MatchStr    string        `yaml:"matchStr"`
	MatchRegex  utils.MyRegex `yaml:"matchRegex"`
	Revert      bool          `yaml:"revert"` // if is true, get result by !
}

func (t MatchField) GetName() string {
	return t.Name
}
func (t MatchField) IsMatch(src map[string]string) bool {

	val, ok := src[t.Key]
	if !ok {
		return false
	}

	if t.MatchStr != "" {
		if t.Revert {
			return !(t.MatchStr == val)
		}
		return t.MatchStr == val
	}

	glog.V(5).Infof(" key %s regex %s result %v ", val, t.MatchRegex.RegexStr, t.MatchRegex.IsMatch(val))

	if t.Revert {
		return !t.MatchRegex.IsMatch(val)
	}
	return t.MatchRegex.IsMatch(val)
}
