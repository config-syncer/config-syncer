/*
Copyright The Kubed Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"fmt"
	"net"
	"os"
	"time"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/pkg/cmds/server"
	srvr "github.com/appscode/kubed/pkg/server"

	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
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

func (f *Framework) RunKubed(stopCh chan struct{}, clusterConfig api.ClusterConfig) error {
	clusterConfig.Save(KubedTestConfigFileDir)

	kubedServer, err := f.NewKubedServer()
	if err != nil {
		return err
	}

	go kubedServer.GenericAPIServer.PrepareRun().Run(stopCh)
	go kubedServer.Operator.Run(stopCh)

	return nil
}

func (f *Framework) NewKubedServer() (*srvr.KubedServer, error) {
	kubedOptions := f.NewKubedOptions()
	config, err := kubedOptions.Config()
	if err != nil {
		return nil, err
	}

	config.OperatorConfig.OperatorNamespace = f.namespace
	config.OperatorConfig.ConfigPath = KubedTestConfigFileDir
	config.OperatorConfig.Test = true

	s, err := config.Complete().New()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (f *Framework) NewKubedOptions() *server.KubedOptions {
	return &server.KubedOptions{
		RecommendedOptions: f.NewKubedServerOptions(),
		OperatorOptions:    server.NewOperatorOptions(),
		StdOut:             os.Stdout,
		StdErr:             os.Stderr,
	}
}

func (f *Framework) NewKubedServerOptions() *options.RecommendedOptions {
	return &options.RecommendedOptions{
		Authentication: &options.DelegatingAuthenticationOptions{
			RemoteKubeConfigFile: f.KubeConfigPath,
			//SkipInClusterLookup:  true,
		},
		Authorization: &options.DelegatingAuthorizationOptions{
			RemoteKubeConfigFile: f.KubeConfigPath,
		},
		CoreAPI: &options.CoreAPIOptions{
			CoreAPIKubeconfigPath: f.KubeConfigPath,
		},
		SecureServing: &options.SecureServingOptionsWithLoopback{
			SecureServingOptions: &options.SecureServingOptions{
				BindPort:    8443,
				BindAddress: net.ParseIP("127.0.0.1"),
			},
		},
		ExtraAdmissionInitializers: func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) { return nil, nil },
		Etcd:                       nil,
		Admission:                  nil,
	}
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
			Group:                 api.GroupName,
			GroupPriorityMinimum:  MaxRetry,
			VersionPriority:       15,
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

const (
	OperatorName      = "kubed-operator"
	OperatorNamespace = "kube-system"
	ContainerOperator = "operator"
	OperatorConfig    = "kubed-config"
)

func (f *Invocation) RestartKubedOperator(config *api.ClusterConfig) error {
	meta := metav1.ObjectMeta{
		Name:      OperatorConfig,
		Namespace: OperatorNamespace,
	}

	err := f.DeleteSecret(meta)
	if err != nil && !kerr.IsNotFound(err) {
		return err
	}

	err = f.WaitUntilSecretDeleted(meta)
	if err != nil {
		return err
	}

	kubeConfig, err := f.KubeConfigSecret(config, meta)
	if err != nil {
		return err
	}

	_, err = f.CreateSecret(kubeConfig)
	if err != nil {
		return err
	}

	pods, err := f.KubeClient.CoreV1().Pods(OperatorNamespace).List(metav1.ListOptions{LabelSelector: "app=kubed"})
	for _, pod := range pods.Items {
		for _, c := range pod.Spec.Containers {
			if c.Name == ContainerOperator {
				err = f.KubeClient.CoreV1().Pods(OperatorNamespace).Delete(pod.Name, deleteInBackground())
				if err != nil {
					return err
				}
				err = f.WaitUntilPodTerminated(pod.ObjectMeta)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	deployment, err := f.KubeClient.AppsV1().Deployments(OperatorNamespace).Get(OperatorName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	return f.WaitUntilDeploymentReady(deployment.ObjectMeta)
}
