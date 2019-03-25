package config

import (
	"cgolib/matchers"
	"cgolib/parsers"
	"fmt"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
)

type Config struct {
	ParseConfig          ParseConfig    `yaml:"parseConfig"`
	Filters              []FilterConfig `yaml:"filters"` // key for matcher name, value for parsers name
	ExternalParseConfigs []string       `yaml:"externalParseConfigs"`
	ExternalFilters      []string       `yaml:"externalFilters"`
}

func (t *Config) ApplyConfig() error {

	// load external parse config file
	for _, fp := range t.ExternalParseConfigs {
		msg, err := ioutil.ReadFile(fp)
		if err != nil {
			return err
		}
		pc := ParseConfig{}
		err = yaml.UnmarshalStrict(msg, &pc)
		if err != nil {
			return err
		}
		err = pc.ApplyConfig()
		if err != nil {
			return err
		}
		err = t.ParseConfig.MergeConfig(&pc, true)
		if err != nil {
			return err
		}
	}

	// load external filters
	for _, fp := range t.ExternalFilters {
		msg, err := ioutil.ReadFile(fp)
		if err != nil {
			return err
		}
		fc := make([]FilterConfig, 0)
		err = yaml.UnmarshalStrict(msg, &fc)
		if err != nil {
			return err
		}

		t.Filters = append(t.Filters, fc...)
	}

	return nil
}

type FilterConfig struct {
	Name     string   `yaml:"name"`
	Matchers []string `yaml:"matchers"`
	Filters  []string `yaml:"filters"`
}

type ParseConfig struct {
	matchers.MatcherConfig `yaml:"match"`
	parsers.ParserConfig   `yaml:"parser"`
	Matchers               map[string]matchers.Matcher
	Parsers                map[string]parsers.Parser
}

func (t *ParseConfig) MergeConfig(pc *ParseConfig, conflict bool) error {
	for k, v := range pc.Matchers {
		_, ok := t.Matchers[k]
		if ok && conflict {
			return fmt.Errorf(" key %s is conflict. ", k)
		}
		if ok {
			glog.Warningf(" will rewrite %s. ", k)
		}
		t.Matchers[k] = v
	}

	for k, v := range pc.Parsers {
		_, ok := t.Parsers[k]
		if ok && conflict {
			return fmt.Errorf(" key %s is conflict. ", k)
		}
		if ok {
			glog.Warningf(" will rewrite %s. ", k)
		}
		t.Parsers[k] = v
	}
	return nil
}

func (t *ParseConfig) ApplyConfig() error {

	// convert to map
	t.Matchers = make(map[string]matchers.Matcher, 0)
	t.Parsers = make(map[string]parsers.Parser, 0)

	// use reflect to load all field
	mc := reflect.ValueOf(t.MatcherConfig)
	for i := 0; i < mc.NumField(); i++ {
		switch mc.Field(i).Type().Kind() {
		case reflect.Slice:
			s := mc.Field(i)
			for j := 0; j < s.Len(); j++ {
				ms := s.Index(j).Interface().(matchers.Matcher)
				t.Matchers[ms.GetName()] = ms
			}
		case reflect.Struct, reflect.Ptr:
			ms := mc.Field(i).Interface().(matchers.Matcher)
			t.Matchers[ms.GetName()] = ms
		default:
			glog.Error("invalid type ", mc.Field(i).Type().Kind())
		}
	}

	// use reflect to load all field
	pc := reflect.ValueOf(t.ParserConfig)
	for i := 0; i < pc.NumField(); i++ {
		switch pc.Field(i).Type().Kind() {
		case reflect.Slice:
			s := pc.Field(i)
			for j := 0; j < s.Len(); j++ {
				ps := s.Index(j).Interface().(parsers.Parser)
				t.Parsers[ps.GetName()] = ps
			}
		case reflect.Struct, reflect.Ptr:
			ps := pc.Field(i).Interface().(parsers.Parser)
			t.Parsers[ps.GetName()] = ps
		default:
			glog.Error("invalid type ", pc.Field(i).Type().Kind())
		}
	}

	return nil
}
