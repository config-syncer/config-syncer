package watcher

import (
	"github.com/appscode/kubed/pkg/events"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

func (w *Controller) watchService() {
	lw := cache.NewListWatchFromClient(
		w.KubeClient.CoreV1().RESTClient(),
		events.Service.String(),
		metav1.NamespaceAll,
		fields.Everything())
	_, controller := cache.NewInformer(lw, &apiv1.Service{}, w.SyncPeriod, eventHandlerFuncs(w))
	go controller.Run(wait.NeverStop)
}
