package utils

import "fmt"

type Field struct {
	Name string
	Key  string
}

func (t Field) GetName() string {
	return t.Name
}

func MergeMap(src, des map[string]string, conflict bool) error {
	for k, v := range des {
		_, ok := src[k]
		if ok && conflict {
			return fmt.Errorf("key is conflict %s. ", k)
		}
		src[k] = v
	}
	return nil
}

func MergeMapWithPrefix(src, des map[string]string, prefix string, conflict bool) error {
	for k, v := range des {
		_, ok := src[prefix+k]
		if ok && conflict {
			return fmt.Errorf("key is conflict %s. ", prefix+k)
		}
		src[prefix+k] = v
	}
	return nil
}
