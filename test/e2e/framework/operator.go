package framework

import (
	"fmt"
	"net"
	"os"
	"time"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/pkg/cmds/server"
	"github.com/appscode/kubed/pkg/operator"
	srvr "github.com/appscode/kubed/pkg/server"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apiserver/pkg/admission"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	apireg "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
)

var (
	svc        *core.Service
	apiService *apireg.APIService
	endpoints  *core.Endpoints
)

func (f *Framework) NewTestKubedOptions(kubeConfigPath string) *server.KubedOptions {
	return &server.KubedOptions{
		RecommendedOptions: f.NewTestKubedServerOptions(kubeConfigPath),
		OperatorOptions:    f.NewTestKubedOperatorOptions(),
		StdOut:             os.Stdout,
		StdErr:             os.Stderr,
	}
}

func (f *Framework) NewTestKubedServerOptions(kubeConfigPath string) *options.RecommendedOptions {
	return &options.RecommendedOptions{
		Authentication: &options.DelegatingAuthenticationOptions{
			RemoteKubeConfigFile: kubeConfigPath,
			//SkipInClusterLookup:  true,
		},
		Authorization: &options.DelegatingAuthorizationOptions{
			RemoteKubeConfigFile: kubeConfigPath,
		},
		CoreAPI: &options.CoreAPIOptions{
			CoreAPIKubeconfigPath: kubeConfigPath,
		},
		SecureServing: &options.SecureServingOptionsWithLoopback{
			SecureServingOptions: &options.SecureServingOptions{
				BindPort:    8443,
				BindAddress: net.ParseIP("127.0.0.1"),
			},
		},
		ExtraAdmissionInitializers: func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) { return nil, nil },
		Etcd:      nil,
		Admission: nil,
	}
}

func (f *Framework) NewTestKubedOperatorOptions() *server.OperatorOptions {
	opt := server.NewOperatorOptions()
	opt.ConfigPath = KubedTestConfigFileDir
	return opt
}

func (f *Framework) NewTestOperatorConfig() *operator.OperatorConfig {
	ocfg := &operator.OperatorConfig{
		Config:            f.KubedServer.Operator.Config,
		ClientConfig:      f.KubedServer.Operator.ClientConfig,
		KubeClient:        f.KubedServer.Operator.KubeClient,
		VoyagerClient:     f.KubedServer.Operator.VoyagerClient,
		StashClient:       f.KubedServer.Operator.StashClient,
		SearchlightClient: f.KubedServer.Operator.SearchlightClient,
		KubeDBClient:      f.KubedServer.Operator.KubeDBClient,
		PromClient:        f.KubedServer.Operator.PromClient,
	}
	ocfg.Test = true
	ocfg.ConfigPath = KubedTestConfigFileDir
	return ocfg
}

func (f *Framework) RunOperator(stopCh chan struct{}, clusterConfig api.ClusterConfig) error {
	var err error
	operatorConfig := f.NewTestOperatorConfig()
	err = clusterConfig.Save(operatorConfig.ConfigPath)
	if err != nil {
		return err
	}

	operatorConfig.OperatorNamespace = f.Namespace()
	f.KubedServer.Operator, err = operatorConfig.New()
	if err != nil {
		return err
	}

	go f.KubedServer.Operator.Run(stopCh)

	return nil
}

func (f *Framework) NewTestKubedServer(kubeConfigPath string) (*srvr.KubedServer, error) {
	kubedOptions := f.NewTestKubedOptions(kubeConfigPath)

	config, err := kubedOptions.Config()
	if err != nil {
		return nil, err
	}

	s, err := config.Complete().New()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (f *Invocation) RegisterAPIService() (err error) {
	svc = f.ServiceForAPIService()
	endpoints = f.LocalEndPoints(svc)
	apiService = f.APIService(svc)

	_, err = f.KubeClient.CoreV1().Services(svc.Namespace).Create(svc)
	if err != nil {
		return
	}

	_, err = f.KubeClient.CoreV1().Endpoints(endpoints.Namespace).Create(endpoints)
	if err != nil {
		return
	}

	_, err = f.KAClient.ApiregistrationV1beta1().APIServices().Create(apiService)

	return
}

func (f *Invocation) DeleteAPIService() {
	f.KAClient.ApiregistrationV1beta1().APIServices().Delete(apiService.Name, deleteInBackground())
	f.KubeClient.CoreV1().Services(svc.Namespace).Delete(svc.Name, deleteInBackground())
	f.KubeClient.CoreV1().Endpoints(endpoints.Namespace).Delete(endpoints.Name, deleteInBackground())
}

func (f *Invocation) APIService(service *core.Service) *apireg.APIService {
	return &apireg.APIService{
		ObjectMeta: metav1.ObjectMeta{
			Name: "v1alpha1.kubed.appscode.com",
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: apireg.APIServiceSpec{
			InsecureSkipTLSVerify: true,
			Group:                api.GroupName,
			GroupPriorityMinimum: MaxRetry,
			VersionPriority:      15,
			Service: &apireg.ServiceReference{
				Name:      service.Name,
				Namespace: service.Namespace,
			},
			Version: api.SchemeGroupVersion.Version,
		},
	}
}

func (f *Invocation) ServiceForAPIService() *core.Service {
	return &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.app,
			Namespace: f.Namespace(),
		},
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{
					Name:       "api",
					Port:       443,
					Protocol:   core.ProtocolTCP,
					TargetPort: intstr.FromInt(8443),
				},
			},
			Type: core.ServiceTypeClusterIP,
		},
	}
}

func (f *Invocation) LocalEndPoints(service *core.Service) *core.Endpoints {
	return &core.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
		},
		Subsets: []core.EndpointSubset{
			{
				Addresses: []core.EndpointAddress{
					{
						IP: "10.0.2.2",
					},
				},
				Ports: []core.EndpointPort{
					{
						Name:     "api",
						Port:     8443,
						Protocol: core.ProtocolTCP,
					},
				},
			},
		},
	}
}

func (f *Framework) EventuallyAPIServerReady() GomegaAsyncAssertion {
	return Eventually(
		func() error {
			apiservice, err := f.KAClient.ApiregistrationV1beta1().APIServices().Get("v1alpha1.kubed.appscode.com", metav1.GetOptions{})
			if err != nil {
				return err
			}
			for _, cond := range apiservice.Status.Conditions {
				if cond.Type == apireg.Available && cond.Status == apireg.ConditionTrue && cond.Reason == "Passed" {
					return nil
				}
			}
			return fmt.Errorf("ApiService not ready yet")
		},
		time.Minute*5,
		time.Microsecond*10,
	)
}

func (f *Framework) DeleteClusterRole(meta metav1.ObjectMeta) error {
	return f.KubeClient.RbacV1().ClusterRoles().Delete(meta.Name, deleteInBackground())
}

func (f *Framework) DeleteClusterRoleBinding(meta metav1.ObjectMeta) error {
	return f.KubeClient.RbacV1().ClusterRoleBindings().Delete(meta.Name, deleteInBackground())
}
