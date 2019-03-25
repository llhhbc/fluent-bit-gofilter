package main

import "C"
import (
	"cgolib/filters"
	"cgolib/modules"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"strings"
)

//export Golib_init
func Golib_init(name, value []string) int {

	src := loadCallIn(name, value)

	args, ok := src["lib_args"]
	if ok {
		err := flag.CommandLine.Parse(strings.Split(args, " "))
		if err != nil {
			fmt.Println("parse lib_args fail ", err)
			return -1
		}
	}

	err := modules.Init(src)
	if err != nil {
		glog.Error("init module fail: ", err)
		return -1
	}

	err = filters.InitPlugins(src)
	if err != nil {
		glog.Error("init plugin fail: ", err)
		return -1
	}

	return 0
}

//export Golib_filter
func Golib_filter(srcName, srcValue []string) int {
	// go will return the result in the slice which called in. so can't append item that max slice's cap.
	// must use cgoAppend instead of append if you want append value.
	// use cgoSetSlice if you want change value of slice or just append.

	src := loadCallIn(srcName, srcValue)

	err := filters.FilterPlugins(src)
	if err != nil {
		glog.Error("do filter fail ", err)
		return -1
	}

	return unLoadCallIn(src, srcName, srcValue)
}

//export Golib_exit
func Golib_exit() int {

	glog.Infoln("go exit")

	err := modules.Exit()
	if err != nil {
		glog.Error("exit module fail ", err)
		return -1
	}

	err = filters.ExitPlugins()
	if err != nil {
		glog.Error("exit filter fail ", err)
		return -1
	}

	return 0
}

func loadCallIn(name, value []string) map[string]string {
	res := make(map[string]string)
	for idx, n := range name {
		res[n] = value[idx]
	}
	return res
}

func unLoadCallIn(src map[string]string, name, value []string) int {
	index := 0
	for k, v := range src {
		name = cgoSetSlice(name, index, k)
		value = cgoSetSlice(value, index, v)
		index++
	}
	glog.V(5).Infoln("unload ok: ", name, value)
	return index
}

// can't append parameters that will over slice's cap.
// if over slice's cap, then go will malloc new memory, the c can't get the results.
func cgoAppend(src []string, parameters ...string) []string {
	if cap(src)-len(src) < len(parameters) {
		glog.Error(" cann't set parameters. slice cap is full. ", len(src), cap(src), len(parameters))
		return src
	}
	src = append(src, parameters...)
	return src
}

func cgoSetSlice(src []string, index int, value string) []string {
	if index > cap(src) {
		glog.Error(" cann't set parameters. index overflow. ", index, cap(src))
		return src
	}

	if index >= len(src) {
		src = append(src, value)
	} else {
		src[index] = value
	}

	return src
}

func main() {}
