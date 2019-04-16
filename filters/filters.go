package filters

import (
	"fmt"
	"github.com/golang/glog"
	"strings"
)

var (
	// filters that register by import go pkg
	registerFilters = make(map[string]Filter)

	// filters that init by config
	initFilters = make(map[string]Filter)
)

func init() {
	registerFilters["parseFilter"] = &ParseFilter{}
}

type Filter interface {
	Init(map[string]string) error
	Filter(map[string]string) error
	Exit() error
}

/*
filter config

// load filter by set filter name:
filters_1  name1
filters_2  name2
filters_3  name3
*/
func InitPlugins(config map[string]string) error {
	var err error

	for k, v := range config {
		if strings.HasPrefix(k, "filters") {
			p, ok := registerFilters[v]
			if !ok {
				return fmt.Errorf("unknow plugin %s. ", v)
			}
			err = p.Init(config)
			if err != nil {
				return err
			}
			initFilters[v] = p
			glog.Infof("init golib plugin %s ok.\n", v)
		}
	}

	return nil
}

func FilterPlugins(src map[string]string) error {
	var err error

	for k, p := range initFilters {
		err = p.Filter(src)
		if err != nil {
			return fmt.Errorf("do filter %s fail %s. ", k, err.Error())
		}
	}

	return nil
}

func ExitPlugins() error {
	var err error

	for _, p := range initFilters {
		err = p.Exit()
		if err != nil {
			return err
		}
	}

	return nil
}
