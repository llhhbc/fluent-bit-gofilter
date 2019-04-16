package main

import (
	"C"
	"cgolib/filters"
	"cgolib/modules"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"strings"
	"time"
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
		fmt.Println("init golib with args ", src["lib_args"])
	}

	err := modules.Init(src)
	if err != nil {
		fmt.Println("init module fail: ", err)
		return -1
	}

	err = filters.InitPlugins(src)
	if err != nil {
		fmt.Println("init plugin fail: ", err)
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
		fmt.Println("do filter fail ", err)
		return -1
	}

	return unLoadCallIn(src, srcName, srcValue)
}

//export Golib_exit
func Golib_exit() int {

	glog.Infoln("go exit")

	defer func() {
		time.Sleep(3 * time.Second) // go need time to free routines. without this , 'fatal: morestack on g0' will happened.
	}()

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
		if n == "log" {
			res[n] = res[n] + value[idx] // log key combine
		} else {
			res[n] = value[idx]
		}
	}
	glog.V(5).Infoln("load ok: ", name, value, res)
	return res
}

func unLoadCallIn(src map[string]string, name, value []string) int {
	index := 0
	for k, v := range src {
		if k == "" {
			// skip  empty name
			continue
		}
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

func main() {
	flag.Parse()
}
