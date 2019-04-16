package modules

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
	"time"
)

type module struct {
	hasInit bool
}

func (t *module) InitModule() {
	t.hasInit = true
}

func (t *module) HasInit() bool {
	return t.hasInit
}

type K8s struct {
	module
	Rest      *rest.Config
	ClientSet kubernetes.Interface
	stopChan  chan struct{}

	PodUid  sync.Map // find by uid
	PodName sync.Map // find by namespace + name
}

var k8s K8s

func GetK8s() *K8s {
	return &k8s
}

/*
parameters:
kube_config          kube client config file
kube_parse_config    kube parse config file --- yaml file
*/
func (t *K8s) Init(src map[string]string) error {
	var err error

	t.Rest, err = clientcmd.BuildConfigFromFlags("", src["kube_config"])
	if err != nil {
		return fmt.Errorf("load kube config %s fail %s. ", src["kube_config"], err.Error())
	}

	t.ClientSet, err = kubernetes.NewForConfig(t.Rest)
	if err != nil {
		return fmt.Errorf("get clientset fail %s. ", err.Error())
	}

	t.stopChan = make(chan struct{}, 0)

	factory := informers.NewSharedInformerFactoryWithOptions(t.ClientSet, 10*time.Second)

	// register informer
	factory.Core().V1().Pods().Informer().AddEventHandlerWithResyncPeriod(cache.ResourceEventHandlerFuncs{
		AddFunc:    t.addPodMap,
		UpdateFunc: t.updatePodMap,
		DeleteFunc: t.deletePodMap,
	}, 10*time.Minute)

	go factory.Start(t.stopChan)
	time.Sleep(time.Second * 3)

	glog.Info("WaitForCacheSync ")
	res := factory.WaitForCacheSync(t.stopChan)
	glog.Infof(" WaitForCacheSync end %v. ", res)

	// load all pod first
	podList, err := factory.Core().V1().Pods().Lister().List(labels.NewSelector())
	for _, pod := range podList {
		t.addPodMap(pod)
	}

	t.InitModule()

	return nil
}

func (t *K8s) Exit() error {
	close(t.stopChan)

	return nil
}

func (t *K8s) addPodMap(obj interface{}) {

	switch d := obj.(type) {
	case *v1.Pod:
		t.PodUid.Store(string(d.UID), d)
		t.PodName.Store(d.Namespace+d.Name, d)
	}
}

func (t *K8s) updatePodMap(oldobj, obj interface{}) {

	switch d := obj.(type) {
	case *v1.Pod:
		t.PodUid.Store(string(d.UID), d)
		t.PodName.Store(d.Namespace+d.Name, d)
	}
}

func (t *K8s) deletePodMap(obj interface{}) {

	switch d := obj.(type) {
	case *v1.Pod:
		t.PodUid.Delete(string(d.UID))
		t.PodName.Delete(d.Namespace + d.Name)
	}
}
