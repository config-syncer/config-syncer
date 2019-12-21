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
	"path/filepath"
	"sync"
	"time"

	"github.com/appscode/go/crypto/rand"

	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"gomodules.xyz/cert/certstore"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	MaxRetry = 200
	NoRetry  = 1

	DefaultEventuallyTimeout         = 5 * time.Minute
	DefaultEventuallyPollingInterval = 2 * time.Second
)

type Framework struct {
	KubeClient         clientset.Interface
	namespace          string
	Mutex              sync.Mutex
	CertStore          *certstore.CertStore
	KubeConfigPath     string
	SelfHostedOperator bool
	ClientConfig       *rest.Config
}

func New(config *rest.Config) *Framework {
	store, err := certstore.NewCertStore(afero.NewMemMapFs(), filepath.Join("", "pki"))
	Expect(err).NotTo(HaveOccurred())

	err = store.InitCA()
	Expect(err).NotTo(HaveOccurred())

	return &Framework{
		namespace: rand.WithUniqSuffix("test-kubed"),

		ClientConfig: config,
		KubeClient:   clientset.NewForConfigOrDie(config),
		CertStore:    store,
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

type Invocation struct {
	*Framework
	app string
}

func (f *Framework) EnsureCreatedCRDs() error {
	_, pErr := f.PromClient.MonitoringV1().Prometheuses(f.namespace).List(metav1.ListOptions{})
	_, sErr := f.PromClient.MonitoringV1().ServiceMonitors(f.namespace).List(metav1.ListOptions{})
	if pErr == nil && sErr == nil {
		return nil
	}
	promCrd := &extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: prom.PrometheusName + "." + prom.SchemeGroupVersion.Group,
		},

		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   prom.SchemeGroupVersion.Group,
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
			Name: prom.ServiceMonitorName + "." + prom.SchemeGroupVersion.Group,
		},

		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   prom.SchemeGroupVersion.Group,
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
			Name: prom.AlertmanagerName + "." + prom.SchemeGroupVersion.Group,
		},

		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   prom.SchemeGroupVersion.Group,
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
