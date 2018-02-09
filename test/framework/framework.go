package framework

import (
	"sync"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/kubed/pkg/server"
	sls "github.com/appscode/searchlight/client"
	scs "github.com/appscode/stash/client"
	vcs "github.com/appscode/voyager/client"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	kcs "github.com/kubedb/apimachinery/client"
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
	crdClient     ecs.ApiextensionsV1beta1Interface
	KubedOperator *server.Operator
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

	config, err := clientcmd.BuildConfigFromFlags(testConfigs.Master, testConfigs.KubeConfig)
	Expect(err).NotTo(HaveOccurred())
	promClient, err := prom.NewForConfig(&prom.DefaultCrdKinds, prom.Group, config)
	Expect(err).NotTo(HaveOccurred())
	crdClient, err := ecs.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())

	return &Framework{
		KubeConfig: config,
		KubeClient: clientset.NewForConfigOrDie(config),
		crdClient:  crdClient,
		namespace:  testConfigs.TestNamespace,
		Config:     testConfigs,
		KubedOperator: &server.Operator{
			KubeClient:        clientset.NewForConfigOrDie(config),
			StashClient:       scs.NewForConfigOrDie(config),
			VoyagerClient:     vcs.NewForConfigOrDie(config),
			SearchlightClient: sls.NewForConfigOrDie(config),
			KubeDBClient:      kcs.NewForConfigOrDie(config),
			PromClient:        promClient,
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
	_, err := f.crdClient.CustomResourceDefinitions().Create(promCrd)
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
	_, err = f.crdClient.CustomResourceDefinitions().Create(svcMonitorCrd)
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
	_, err = f.crdClient.CustomResourceDefinitions().Create(alertMgr)
	if err != nil {
		return err
	}
	return nil
}
