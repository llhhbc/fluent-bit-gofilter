package parsers

import (
	"cgolib/utils"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"time"
)

type ParseTimer struct {
	utils.Field `yaml:",inline"`
	TimeFormat  string `yaml:"time_format"`
}

func (t ParseTimer) Parse(src map[string]string) error {

	if src[t.Key] == "" {
		glog.Warningf(" ParseTimer %s key %s is not exists. ", t.Name, t.Key)
		return nil
	}

	tm, err := time.Parse(t.TimeFormat, src[t.Key])
	if err != nil {
		return fmt.Errorf("ParseTimer %s parse time %s fail %s. ", t.Name, src[t.Key], err)
	}

	src[t.Key] = tm.Format(time.RFC3339)

	return nil
}

type ParseRegex struct {
	utils.Field `yaml:",inline"`
	Regex       utils.MyRegex `yaml:"regex"`
	Conflict    bool          `yaml:"conflict"`
	PreserveKey bool          `yaml:"preserve_key"`
}

func (t ParseRegex) Parse(src map[string]string) error {

	if src[t.Key] == "" {
		glog.Warningf(" ParseRegex %s key %s is not exists. ", t.Name, t.Key)
		return nil
	}

	values := t.Regex.Regex.FindStringSubmatch(src[t.Key])
	keys := t.Regex.Regex.SubexpNames()
	if len(keys) != len(values) {
		glog.Warningf("ParseRegex %s key %v values %s of %s. %v ", t.Name, keys, values, src[t.Key], src)
		return nil
	}

	if !t.PreserveKey {
		delete(src, t.Key)
	}

	for i := 0; i < len(keys); i++ {
		_, ok := src[keys[i]]
		if ok && t.Conflict {
			return fmt.Errorf("ParseRegex %s key %s is exists. ", t.Name, keys[i])
		}
		src[keys[i]] = values[i]
	}

	return nil
}

type ParseJson struct {
	utils.Field `yaml:",inline"`
	Conflict    bool `yaml:"conflict"`
	PreserveKey bool `yaml:"preserve_key"`
}

func (t ParseJson) Parse(src map[string]string) error {

	if src[t.Key] == "" {
		glog.Warningf(" parseJson %s key %s is not exists. ", t.Name, t.Key)
		return nil
	}
	glog.V(3).Infoln("parseJson %s key %s name %s. ", t.Name, t.Key)
	glog.V(5).Infoln("parseJson %s key %s name %s src %s. ", t.Name, t.Key, src[t.Key])

	res := make(map[string]string, 0)
	err := json.Unmarshal([]byte(src[t.Key]), &res)
	if err != nil {
		return fmt.Errorf("parseJson %s parse fail %s src %v. ", t.Name, err, src)
	}

	err = utils.MergeMap(src, res, t.Conflict)
	if err != nil {
		return fmt.Errorf("parseJson %s parse fail %s. ", t.Name, err)
	}

	if !t.PreserveKey {
		delete(src, t.Key)
	}

	return nil
}
