package framework

import (
	"sync"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/kubed/pkg/operator"
	sls "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	scs "github.com/appscode/stash/client/typed/stash/v1alpha1"
	vcs "github.com/appscode/voyager/client/typed/voyager/v1beta1"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	kcs "github.com/k8sdb/apimachinery/client/typed/kubedb/v1alpha1"
	. "github.com/onsi/gomega"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	ecs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	MaxRetry = 200
	NoRetry  = 1

	DefaultEventuallyTimeout         = 5 * time.Minute
	DefaultEventuallyPollingInterval = 2 * time.Second
)

type Framework struct {
	KubeConfig    *rest.Config
	KubeClient    clientset.Interface
	KubedOperator *operator.Operator
	Config        E2EConfig
	namespace     string
	Mutex         sync.Mutex
}

type Invocation struct {
	*Framework
	app string
}

func New() *Framework {
	testConfigs.validate()

	c, err := clientcmd.BuildConfigFromFlags(testConfigs.Master, testConfigs.KubeConfig)
	Expect(err).NotTo(HaveOccurred())
	promClient, err := prom.NewForConfig(c)
	Expect(err).NotTo(HaveOccurred())
	crdClient, err := ecs.NewForConfig(c)
	Expect(err).NotTo(HaveOccurred())

	return &Framework{
		KubeConfig: c,
		KubeClient: clientset.NewForConfigOrDie(c),
		namespace:  testConfigs.TestNamespace,
		Config:     testConfigs,
		KubedOperator: &operator.Operator{
			KubeClient:        clientset.NewForConfigOrDie(c),
			StashClient:       scs.NewForConfigOrDie(c),
			VoyagerClient:     vcs.NewForConfigOrDie(c),
			SearchlightClient: sls.NewForConfigOrDie(c),
			KubeDBClient:      kcs.NewForConfigOrDie(c),
			PromClient:        promClient,
			CRDClient:         crdClient,
		},
	}
}

func (f *Framework) Invoke() *Invocation {
	return &Invocation{
		Framework: f,
		app:       rand.WithUniqSuffix("kubed-e2e"),
	}
}

func (f *Invocation) App() string {
	return f.app
}

func (f *Framework) EnsureCreatedCRDs() error {
	_, pErr := f.KubedOperator.PromClient.Prometheuses(f.Config.TestNamespace).List(metav1.ListOptions{})
	_, sErr := f.KubedOperator.PromClient.ServiceMonitors(f.Config.TestNamespace).List(metav1.ListOptions{})
	if pErr == nil && sErr == nil {
		return nil
	}
	promCrd := &extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: prom.PrometheusName + "." + prom.Group,
		},

		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   prom.Group,
			Version: prom.Version,
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural: prom.PrometheusName,
				Kind:   prom.PrometheusesKind,
			},
		},
	}
	_, err := f.KubedOperator.CRDClient.CustomResourceDefinitions().Create(promCrd)
	if err != nil {
		return err
	}
	svcMonitorCrd := &extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: prom.ServiceMonitorName + "." + prom.Group,
		},

		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   prom.Group,
			Version: prom.Version,
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural: prom.ServiceMonitorName,
				Kind:   prom.ServiceMonitorsKind,
			},
		},
	}
	_, err = f.KubedOperator.CRDClient.CustomResourceDefinitions().Create(svcMonitorCrd)
	if err != nil {
		return err
	}

	alertMgr := &extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: prom.AlertmanagerName + "." + prom.Group,
		},

		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   prom.Group,
			Version: prom.Version,
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural: prom.AlertmanagerName,
				Kind:   prom.AlertmanagersKind,
			},
		},
	}
	_, err = f.KubedOperator.CRDClient.CustomResourceDefinitions().Create(alertMgr)
	if err != nil {
		return err
	}
	return nil
}
