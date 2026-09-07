/*
Copyright AppsCode Inc. and Contributors

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

package cluster

import (
	"context"
	"fmt"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	OpenShiftClusterMonitoringNamespace = "openshift-monitoring"
	OpenShiftClusterPrometheus          = "k8s"
	OpenShiftClusterAlertmanager        = "main"
	OpenShiftThanosQuerierService       = "thanos-querier"

	OpenShiftUserWorkloadMonitoringNamespace = "openshift-user-workload-monitoring"
	OpenShiftUserWorkloadPrometheus          = "user-workload"
	OpenShiftUserWorkloadAlertmanager        = "user-workload"
)

func IsOpenShiftManaged(mapper meta.RESTMapper) bool {
	if _, err := mapper.RESTMappings(schema.GroupKind{
		Group: "project.openshift.io",
		Kind:  "Project",
	}); err == nil {
		return true
	}
	return false
}

// GetOpenShiftAppsDomain fetches the default *.apps.<cluster_name>.<base_domain> domain for OpenShift
func GetOpenShiftAppsDomain(kc client.Client) (string, error) {
	var ing unstructured.Unstructured
	ing.SetAPIVersion("operator.openshift.io/v1")
	ing.SetKind("IngressController")
	key := client.ObjectKey{Namespace: "openshift-ingress-operator", Name: "default"}
	if err := kc.Get(context.Background(), key, &ing); err != nil {
		return "", err
	}
	domain, found, err := unstructured.NestedString(ing.Object, "status", "domain")
	if err != nil {
		return "", err
	}
	if !found || domain == "" {
		return "", fmt.Errorf("status.domain not found in IngressController")
	}
	return domain, nil
}

// GetOpenShiftServiceSigner fetches the OpenShift service signer CA certificate
func GetOpenShiftServiceSigner(kc client.Client) ([]byte, error) {
	var cm core.ConfigMap
	err := kc.Get(context.TODO(), client.ObjectKey{Namespace: "kube-public", Name: "openshift-service-ca.crt"}, &cm)
	if err != nil {
		return nil, err
	}
	return []byte(cm.Data["service-ca.crt"]), nil
}
