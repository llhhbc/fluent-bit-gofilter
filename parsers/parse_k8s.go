package parsers

import (
	"cgolib/modules"
	"cgolib/utils"
	"fmt"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
)

type ParseK8sUid struct {
	utils.Field `yaml:",inline"`
	Prefix      string `yaml:"prefix"`
	Conflict    bool   `yaml:"conflict"`
}

const k8sPrefix = "k8s/"

func (t *ParseK8sUid) Parse(src map[string]string) error {

	k8s := modules.GetK8s()
	if !k8s.HasInit() {
		return fmt.Errorf("ParseK8sUid %s need init module k8s. ", t.Name)
	}

	if src[t.Key] == "" {
		glog.Warningf("ParseK8sUid %s key %s is empty. ", t.Name, t.Key)
		return nil
	}

	p, ok := k8s.PodUid.Load(src[t.Key])
	if !ok {
		glog.Errorf("ParseK8sUid %s key %s is not exist. ", t.Name, src[t.Key])
		return nil
	}
	pod, ok := p.(*v1.Pod)
	if !ok {
		glog.Errorf("ParseK8sUid %s key %s convert to pod fail %v. ", t.Name, src[t.Key], p)
		return nil
	}

	if t.Prefix == "" {
		t.Prefix = k8sPrefix
	}
	utils.MergeMap(src, map[string]string{
		"pod_name":       pod.Name,
		"container_name": pod.Name,
		"namespace_name": pod.Namespace,
	}, false)

	err := utils.MergeMapWithPrefix(src, pod.Labels, t.Prefix+"labels/", t.Conflict)
	if err != nil {
		return fmt.Errorf("ParseK8sUid %s key %s merge labels fail %s. ", t.Name, src[t.Key], err)
	}
	err = utils.MergeMapWithPrefix(src, pod.Annotations, t.Prefix+"annotations/", t.Conflict)
	if err != nil {
		return fmt.Errorf("ParseK8sUid %s key %s merge annotation fail %s. ", t.Name, src[t.Key], err)
	}

	return nil
}

const k8sContainerLogName = `(?<tag>[^.]+)?\.?(?<pod_name>[a-z0-9](?:[-a-z0-9]*[a-z0-9])?(?:\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*)_(?<namespace_name>[^_]+)_(?<container_name>.+)-(?<docker_id>[a-z0-9]{64})\.log$`

type ParseK8sName struct {
	utils.Field  `yaml:",inline"`
	NamespaceKey string `yaml:"namespaceKey"`
	Prefix       string `yaml:"prefix"`
	Conflict     bool   `yaml:"conflict"`
}

func (t *ParseK8sName) Parse(src map[string]string) error {
	k8s := modules.GetK8s()
	if !k8s.HasInit() {
		return fmt.Errorf("ParseK8sName %s need init module k8s. ", t.Name)
	}

	if t.Key == "" {
		t.Key = "pod_name"
	}

	if src[t.Key] == "" {
		glog.Warningf("ParseK8sName %s key pod_name is empty %v. ", t.Name, src)
		return nil
	}

	if t.NamespaceKey == "" {
		t.NamespaceKey = "namespace_name"
	}
	if src[t.NamespaceKey] == "" {
		glog.Warningf("ParseK8sName %s key namespace_name is empty %v. ", t.Name, src)
		return nil
	}

	p, ok := k8s.PodName.Load(src[t.NamespaceKey] + src[t.Key])
	if !ok {
		glog.Errorf("ParseK8sName %s key %s is not exist. ", t.Name, src[t.NamespaceKey]+src[t.Key])
		return nil
	}
	pod, ok := p.(*v1.Pod)
	if !ok {
		glog.Errorf("ParseK8sName %s key %s convert to pod fail %v. ", t.Name, src[t.NamespaceKey]+src[t.Key], p)
		return nil
	}

	if t.Prefix == "" {
		t.Prefix = k8sPrefix
	}
	err := utils.MergeMapWithPrefix(src, pod.Labels, t.Prefix+"labels/", t.Conflict)
	if err != nil {
		return fmt.Errorf("ParseK8sUid %s key %s merge labels fail %s. ", t.Name, src[t.Key], err)
	}
	err = utils.MergeMapWithPrefix(src, pod.Annotations, t.Prefix+"annotations/", t.Conflict)
	if err != nil {
		return fmt.Errorf("ParseK8sUid %s key %s merge annotation fail %s. ", t.Name, src[t.Key], err)
	}

	return nil
}
