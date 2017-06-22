package indexers

import (
	"net/http"
	"sync"
	"time"

	"github.com/appscode/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	kcache "k8s.io/client-go/tools/cache"
)

const (
	// Resync period for the kube controller loop.
	resyncPeriod = 5 * time.Minute
)

type ReverseIndexer struct {
	// kubeClient makes calls to API Server and registers calls with API Server
	kubeClient clientset.Interface

	// reverseRecordMap pod to service object.
	reverseRecordMap map[string][]*apiv1.Service

	// cacheLock protecting the cache. caller is responsible for using
	// the cacheLock before invoking methods on cache the cache is not
	// thread-safe, and the caller can guarantee thread safety by using
	// the cacheLock
	cacheLock sync.RWMutex

	// serviceController invokes registered callbacks when services change.
	serviceController kcache.Controller
	// servicesStore that contains all the services in the system.
	servicesStore kcache.Store

	// Initial timeout for endpoints and services to be synced from APIServer
	initialSyncTimeout time.Duration

	apiHandler http.Handler
}

// NeverStop may be passed to Until to make it never stop.
var NeverStop <-chan struct{} = make(chan struct{})

func NewReverseIndexer(client clientset.Interface, timeout time.Duration) *ReverseIndexer {
	ri := &ReverseIndexer{
		kubeClient:         client,
		cacheLock:          sync.RWMutex{},
		reverseRecordMap:   make(map[string][]*apiv1.Service),
		initialSyncTimeout: timeout,
		apiHandler:         &reverseIndexAPIHandlers{},
	}

	ri.setServiceWatcher()

	return ri
}

func (ri *ReverseIndexer) Start() {
	log.Infoln("Starting serviceController")
	go ri.serviceController.Run(NeverStop)

	// Wait synchronously for the initial list operations to be
	// complete of endpoints and services from APIServer.
	ri.waitForResourceSyncedOrDie()
}

func (ri *ReverseIndexer) waitForResourceSyncedOrDie() {
	// Wait for both controllers have completed an initial resource listing
	timeout := time.After(ri.initialSyncTimeout)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-timeout:
			log.Fatalf("Timeout waiting for initialization")
		case <-ticker.C:
			if ri.serviceController.HasSynced() {
				log.Infoln("Initialized services from apiserver")
				return
			}
			log.Infof("Waiting for services and endpoints to be initialized from apiserver...")
		}
	}
}

func (ri *ReverseIndexer) setServiceWatcher() {
	// Returns a cache.ListWatch that gets all changes to services.
	ri.servicesStore, ri.serviceController = kcache.NewInformer(
		&kcache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return ri.kubeClient.CoreV1().Services(apiv1.NamespaceAll).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return ri.kubeClient.CoreV1().Services(apiv1.NamespaceAll).Watch(options)
			},
		},
		&apiv1.Service{},
		resyncPeriod,
		kcache.ResourceEventHandlerFuncs{
			AddFunc:    ri.newService,
			DeleteFunc: ri.removeService,
			UpdateFunc: ri.updateService,
		},
	)
}

func (ri *ReverseIndexer) newService(obj interface{}) {
	if service, ok := assertIsService(obj); ok {
		log.Infof("New service: %v", service.Name)
		log.V(5).Infof("Service details: %v", service)

	}
}

func (ri *ReverseIndexer) removeService(obj interface{}) {
	if _, ok := assertIsService(obj); ok {

	}
}

func (ri *ReverseIndexer) updateService(oldObj, newObj interface{}) {
	if _, ok := assertIsService(newObj); ok {
		if _, ok := assertIsService(oldObj); ok {

		}
	}
}

func assertIsService(obj interface{}) (*apiv1.Service, bool) {
	if service, ok := obj.(*apiv1.Service); ok {
		return service, ok
	} else {
		log.Errorf("Type assertion failed! Expected 'Service', got %T", service)
		return nil, ok
	}
}
