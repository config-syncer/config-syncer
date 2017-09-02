package indexers

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/appscode/go/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
)

func newTestReverseIndexer() *ReverseIndexer {
	c, err := ensureIndex(os.TempDir()+"/index", "indexer")
	if err != nil {
		log.Fatal(err)
	}
	return &ReverseIndexer{
		kubeClient: fake.NewSimpleClientset(
			newPod("foo-pod-1"),
			newPod("foo-pod-2"),
		),
		dataChan: make(chan interface{}, 1),
		index:    c,
	}
}

func TestNewService(t *testing.T) {
	defer os.RemoveAll(os.TempDir() + "/index")

	ri := newTestReverseIndexer()
	ri.dataChan <- newService()
	ri.AddService()

	pod := newPod("foo-pod-1")
	if rawdata, err := ri.index.GetInternal(namespacerKey(pod.ObjectMeta)); err == nil {
		var svc []*v1.Service
		err := json.Unmarshal(rawdata, &svc)
		if err != nil {
			t.Fatal(err)
		}
		if !equalService(svc[0], newService()) {
			t.Errorf("Service did not matched")
		}
	} else {
		t.Errorf("Service did not found in cache")
	}

	pod = newPod("foo-pod-2")
	if rawdata, err := ri.index.GetInternal(namespacerKey(pod.ObjectMeta)); err == nil {
		var svc []*v1.Service
		err := json.Unmarshal(rawdata, &svc)
		if err != nil {
			t.Fatal(err)
		}
		if !equalService(svc[0], newService()) {
			t.Errorf("Service did not matched")
		}
	} else {
		t.Errorf("Service did not found in cache")
	}

	pod = newPod("foo-pod-3")
	if res, err := ri.index.GetInternal(namespacerKey(pod.ObjectMeta)); err == nil {
		if len(res) > 0 {
			t.Errorf("Service Found, expected Not Found")
		}
	}
}

func TestRemoveService(t *testing.T) {
	defer os.RemoveAll(os.TempDir() + "/index")

	ri := newTestReverseIndexer()

	service := newService()
	ri.dataChan <- service
	ri.AddService()
	pod := newPod("foo-pod-1")
	if rawdata, err := ri.index.GetInternal(namespacerKey(pod.ObjectMeta)); err == nil {
		var svc []*v1.Service
		err := json.Unmarshal(rawdata, &svc)
		if err != nil {
			t.Fatal(err)
		}
		if !equalService(svc[0], service) {
			t.Errorf("Service did not matched")
		}
	} else {
		t.Errorf("Service did not found in cache")
	}

	ri.dataChan <- service
	ri.RemoveService()

	pod = newPod("foo-pod-1")
	if res, err := ri.index.GetInternal(namespacerKey(pod.ObjectMeta)); err == nil {
		if len(res) > 0 {
			fmt.Println(string(res))
			t.Errorf("Service Found, expected Not Found")
		}
	}

	pod = newPod("foo-pod-2")
	if res, err := ri.index.GetInternal(namespacerKey(pod.ObjectMeta)); err == nil {
		if len(res) > 0 {
			t.Errorf("Service Found, expected Not Found")
		}
	}
}

func newService() *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"service-name": "foo",
			},
		},
	}
}

func newPod(name string) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels: map[string]string{
				"service-name": "foo",
			},
		},
	}
}
