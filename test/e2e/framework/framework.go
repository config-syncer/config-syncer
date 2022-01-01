/*
Copyright The Config Syncer Authors.

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

	. "github.com/onsi/gomega"
	"gomodules.xyz/blobfs"
	"gomodules.xyz/cert/certstore"
	"gomodules.xyz/x/crypto/rand"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Framework struct {
	KubeClient     clientset.Interface
	namespace      string
	Mutex          sync.Mutex
	CertStore      *certstore.CertStore
	KubeConfigPath string
	ClientConfig   *rest.Config
}

func New(config *rest.Config) *Framework {
	store, err := certstore.New(blobfs.NewInMemoryFS(), filepath.Join("", "pki"))
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

func (fi *Invocation) App() string {
	return fi.app
}

type Invocation struct {
	*Framework
	app string
}
