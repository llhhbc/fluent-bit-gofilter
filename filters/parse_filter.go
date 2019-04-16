package filters

import (
	"cgolib/config"
	"fmt"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ParseFilter struct {
	Config *config.Config
}

const (
	MatchOne = "matchOne"
	MatchAll = "matchAll"
)

/*
parse_config_file
*/
func (t *ParseFilter) Init(src map[string]string) error {

	file, ok := src["parse_config_file"]
	if !ok {
		return fmt.Errorf("no parse_config_file has config. ")
	}

	msg, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.UnmarshalStrict(msg, &t.Config)
	if err != nil {
		return err
	}

	err = t.Config.ApplyConfig()
	if err != nil {
		return err
	}

	res, _ := yaml.Marshal(t.Config)
	glog.Infof("load config ok: %s ", res)

	return nil
}
func (t ParseFilter) Filter(src map[string]string) error {

	if t.Config == nil {
		return nil
	}

	var isMatch bool
	for idx, fc := range t.Config.Filters {
		glog.V(3).Infof("begin do filter %s %s .", fc.Name, fc.MatchType)
		isMatch = false
		for _, m := range fc.Matchers {
			mm, ok := t.Config.ParseConfig.Matchers[m]
			if !ok {
				glog.Errorf("in filter %s matcher %s not found. ", fc.Name, m)
				continue
			}
			if mm.IsMatch(src) {
				isMatch = true
				if fc.MatchType == MatchOne {
					break
				}
			} else {
				isMatch = false
				if fc.MatchType == MatchAll {
					break
				}
			}
		}
		if !isMatch {
			glog.V(1).Infof(" filter %s not match. ", fc.Name)
			continue
		}
		// do parsers
		for _, p := range fc.Parsers {
			pp, ok := t.Config.ParseConfig.Parsers[p]
			if !ok {
				glog.Errorf("in filter %s parser %s not found. ", fc.Name, p)
				continue
			}
			err := pp.Parse(src)
			if err != nil {
				glog.Errorf("in filter %s idx %d parser %s fail %s. ", fc.Name, idx, pp.GetName(), err)
				continue
			}
			glog.V(3).Infof("in filter %s idx %d parser %s ok. ", fc.Name, idx, pp.GetName())
		}
	}

	return nil
}
func (t ParseFilter) Exit() error {
	return nil
}
