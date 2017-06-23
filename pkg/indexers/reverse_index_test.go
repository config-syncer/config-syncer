package indexers

import (
	"fmt"
	"sync"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/v1"
)

func newTestReverseIndexer() *ReverseIndexer {
	return &ReverseIndexer{
		kubeClient: fake.NewSimpleClientset(
			newPod("foo-pod-1"),
			newPod("foo-pod-2"),
		),
		cacheLock:          sync.RWMutex{},
		reverseRecordMap:   make(map[string][]*v1.Service),
		initialSyncTimeout: time.Minute * 5,
	}
}

func TestNewService(t *testing.T) {
	ri := newTestReverseIndexer()
	ri.newService(newService())

	pod := newPod("foo-pod-1")
	if svc, ok := ri.reverseRecordMap[namespacerKey(pod.ObjectMeta)]; ok {
		if !equalService(svc[0], newService()) {
			t.Errorf("Service did not matched")
		}
	} else {
		t.Errorf("Service did not found in cache")
	}

	pod = newPod("foo-pod-2")
	if svc, ok := ri.reverseRecordMap[namespacerKey(pod.ObjectMeta)]; ok {
		if !equalService(svc[0], newService()) {
			t.Errorf("Service did not matched")
		}
	} else {
		t.Errorf("Service did not found in cache")
	}

	pod = newPod("foo-pod-3")
	if _, ok := ri.reverseRecordMap[namespacerKey(pod.ObjectMeta)]; ok {
		t.Errorf("Service Found, expected Not Found")
	}
}

func TestRemoveService(t *testing.T) {
	ri := newTestReverseIndexer()

	service := newService()
	ri.newService(service)
	pod := newPod("foo-pod-1")
	if svc, ok := ri.reverseRecordMap[namespacerKey(pod.ObjectMeta)]; ok {
		if !equalService(svc[0], service) {
			t.Errorf("Service did not matched")
		}
	} else {
		t.Errorf("Service did not found in cache")
	}

	ri.removeService(service)

	fmt.Println(ri.reverseRecordMap)

	pod = newPod("foo-pod-1")
	if _, ok := ri.reverseRecordMap[namespacerKey(pod.ObjectMeta)]; ok {
		t.Errorf("Service Found, expected Not Found")
	}

	pod = newPod("foo-pod-2")
	if _, ok := ri.reverseRecordMap[namespacerKey(pod.ObjectMeta)]; ok {
		t.Errorf("Service Found, expected Not Found")
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
