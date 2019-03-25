package modules

import (
	"fmt"
	"github.com/golang/glog"
	"strings"
)

var (
	registerModel = make(map[string]Module, 0)

	initModel = make(map[string]Module, 0)
)

type Module interface {
	Init(src map[string]string) error
	HasInit() bool
	Exit() error
}

func init() {
	registerModel["k8s"] = GetK8s()
}

/*
init_modules   moules split by ','
*/
func Init(src map[string]string) error {

	initModules := src["init_modules"]

	if initModules == "" {
		glog.Warningln("no init_modules init.")
		return nil
	}

	modules := strings.Split(initModules, ",")
	for _, m := range modules {
		r, ok := registerModel[m]
		if !ok {
			return fmt.Errorf("unknow module %s. ", m)
		}
		err := r.Init(src)
		if err != nil {
			return fmt.Errorf("init module %s fail %s. ", m, err.Error())
		}
		initModel[m] = r
	}

	return nil
}

func Exit() error {

	for k, v := range initModel {
		err := v.Exit()
		if err != nil {
			glog.Errorf("exit module %s fail %s. ", k, err.Error())
		}
	}

	return nil
}
