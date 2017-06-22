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
	"k8s.io/client-go/pkg/util/wait"
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
	go ri.serviceController.Run(wait.NeverStop)

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
		// AddFunc:    ri.newService,
		// DeleteFunc: ri.removeService,
		// UpdateFunc: ri.updateService,
		},
	)
}
