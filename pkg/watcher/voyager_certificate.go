package watcher

import (
	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	tapi "github.com/appscode/voyager/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (c *Controller) WatchVoyagerCertificates() {
	if !util.IsPreferredAPIResource(c.KubeClient, tapi.V1beta1SchemeGroupVersion.String(), tapi.ResourceKindCertificate) {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", tapi.V1beta1SchemeGroupVersion.String(), tapi.ResourceKindCertificate)
		return
	}
	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return c.VoyagerClient.Certificates(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return c.VoyagerClient.Certificates(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&tapi.Certificate{},
		c.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: func(obj interface{}) {
				if cert, ok := obj.(*tapi.Certificate); ok {
					log.Infof("Certificate %s@%s deleted", cert.Name, cert.Namespace)
					c.Saver.Save(cert.ObjectMeta, obj)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
