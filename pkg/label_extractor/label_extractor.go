package label_extractor

import (
	"sync"

	"github.com/hashicorp/golang-lru"
	"k8s.io/client-go/kubernetes"
)

type ExtractDockerLabel struct {
	kubeClient kubernetes.Interface

	enable    bool
	twoQCache *lru.TwoQueueCache
	lock      sync.RWMutex
}

type RegistrySecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func New(kubeClient kubernetes.Interface) *ExtractDockerLabel {
	return &ExtractDockerLabel{
		kubeClient: kubeClient,
	}
}

func (l *ExtractDockerLabel) Configure(enable bool) {
	l.lock.Lock()
	l.enable = enable
	l.twoQCache, _ = lru.New2Q(128)
	l.lock.Unlock()
}
