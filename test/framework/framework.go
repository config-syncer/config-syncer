package framework

import (
	"sync"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/kubed/pkg/operator"
	sls "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	scs "github.com/appscode/stash/client/typed/stash/v1alpha1"
	vcs "github.com/appscode/voyager/client/typed/voyager/v1beta1"
	kcs "github.com/k8sdb/apimachinery/client/typed/kubedb/v1alpha1"
	. "github.com/onsi/gomega"
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
