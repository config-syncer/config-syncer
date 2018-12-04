package operator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/appscode/go/log"
	apiext_util "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/appscode/kutil/discovery"
	searchlight_api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	stash_api "github.com/appscode/stash/apis/stash/v1alpha1"
	voyager_api "github.com/appscode/voyager/apis/voyager/v1beta1"
	kubedb_api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	apps "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v1"
	batch_v1 "k8s.io/api/batch/v1"
	batch_v1beta1 "k8s.io/api/batch/v1beta1"
	certificates "k8s.io/api/certificates/v1beta1"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	networking "k8s.io/api/networking/v1"
	policy "k8s.io/api/policy/v1beta1"
	rbac "k8s.io/api/rbac/v1"
	storage_v1 "k8s.io/api/storage/v1"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir" // admission "k8s.io/api/admission/v1beta1"
)

func TestRestMapper(t *testing.T) {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube/config")

	_, err := os.Stat(kubeconfigPath)
	if err != nil { //kubeconfig file not found. so skip testing.
		return
	}

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	kc := kubernetes.NewForConfigOrDie(config)
	var crdClient crd_cs.ApiextensionsV1beta1Interface
	crdClient = crd_cs.NewForConfigOrDie(config)

	crds := []*crd_api.CustomResourceDefinition{
		// voyager
		voyager_api.Ingress{}.CustomResourceDefinition(),
		voyager_api.Certificate{}.CustomResourceDefinition(),
		// stash
		stash_api.Restic{}.CustomResourceDefinition(),
		stash_api.Recovery{}.CustomResourceDefinition(),
		// searchlight
		searchlight_api.ClusterAlert{}.CustomResourceDefinition(),
		searchlight_api.NodeAlert{}.CustomResourceDefinition(),
		searchlight_api.PodAlert{}.CustomResourceDefinition(),

		// kubedb
		kubedb_api.Postgres{}.CustomResourceDefinition(),
		kubedb_api.Elasticsearch{}.CustomResourceDefinition(),
		kubedb_api.MySQL{}.CustomResourceDefinition(),
		kubedb_api.MongoDB{}.CustomResourceDefinition(),
		kubedb_api.Redis{}.CustomResourceDefinition(),
		kubedb_api.Memcached{}.CustomResourceDefinition(),
		kubedb_api.Snapshot{}.CustomResourceDefinition(),
		kubedb_api.DormantDatabase{}.CustomResourceDefinition(),
	}
	apiext_util.RegisterCRDs(crdClient, crds)

	restmapper, err := discovery.LoadRestMapper(kc.Discovery())
	if err != nil {
		t.Fatal(err)
	}

	data := []struct {
		in  interface{}
		out schema.GroupVersionResource
	}{
		{&apps.ControllerRevision{}, apps.SchemeGroupVersion.WithResource("controllerrevisions")},
		{&apps.Deployment{}, apps.SchemeGroupVersion.WithResource("deployments")},
		{&apps.ReplicaSet{}, apps.SchemeGroupVersion.WithResource("replicasets")},
		{&apps.StatefulSet{}, apps.SchemeGroupVersion.WithResource("statefulsets")},
		{&autoscaling.HorizontalPodAutoscaler{}, autoscaling.SchemeGroupVersion.WithResource("horizontalpodautoscalers")},
		{&batch_v1.Job{}, batch_v1.SchemeGroupVersion.WithResource("jobs")},
		{&batch_v1beta1.CronJob{}, batch_v1beta1.SchemeGroupVersion.WithResource("cronjobs")},
		{&certificates.CertificateSigningRequest{}, certificates.SchemeGroupVersion.WithResource("certificatesigningrequests")},
		{&core.ComponentStatus{}, core.SchemeGroupVersion.WithResource("componentstatuses")},
		{&core.ConfigMap{}, core.SchemeGroupVersion.WithResource("configmaps")},
		{&core.Endpoints{}, core.SchemeGroupVersion.WithResource("endpoints")},
		{&core.Event{}, core.SchemeGroupVersion.WithResource("events")},
		{&core.LimitRange{}, core.SchemeGroupVersion.WithResource("limitranges")},
		{&core.Namespace{}, core.SchemeGroupVersion.WithResource("namespaces")},
		{&core.Node{}, core.SchemeGroupVersion.WithResource("nodes")},
		{&core.PersistentVolumeClaim{}, core.SchemeGroupVersion.WithResource("persistentvolumeclaims")},
		{&core.PersistentVolume{}, core.SchemeGroupVersion.WithResource("persistentvolumes")},
		{&core.PodTemplate{}, core.SchemeGroupVersion.WithResource("podtemplates")},
		{&core.Pod{}, core.SchemeGroupVersion.WithResource("pods")},
		{&core.ReplicationController{}, core.SchemeGroupVersion.WithResource("replicationcontrollers")},
		{&core.ResourceQuota{}, core.SchemeGroupVersion.WithResource("resourcequotas")},
		{&core.Secret{}, core.SchemeGroupVersion.WithResource("secrets")},
		{&core.ServiceAccount{}, core.SchemeGroupVersion.WithResource("serviceaccounts")},
		{&core.Service{}, core.SchemeGroupVersion.WithResource("services")},
		{&apps.DaemonSet{}, apps.SchemeGroupVersion.WithResource("daemonsets")},
		{&extensions.Ingress{}, extensions.SchemeGroupVersion.WithResource("ingresses")},
		{&networking.NetworkPolicy{}, networking.SchemeGroupVersion.WithResource("networkpolicies")},
		{&policy.PodDisruptionBudget{}, policy.SchemeGroupVersion.WithResource("poddisruptionbudgets")},
		{&rbac.ClusterRoleBinding{}, rbac.SchemeGroupVersion.WithResource("clusterrolebindings")},
		{&rbac.ClusterRole{}, rbac.SchemeGroupVersion.WithResource("clusterroles")},
		{&rbac.RoleBinding{}, rbac.SchemeGroupVersion.WithResource("rolebindings")},
		{&rbac.Role{}, rbac.SchemeGroupVersion.WithResource("roles")},
		//{&settings.PodPreset{}, settings.SchemeGroupVersion.WithResource("podpresets")},
		{&storage_v1.StorageClass{}, storage_v1.SchemeGroupVersion.WithResource("storageclasses")},

		// voyager
		//{&voyager_api.Ingress{}, voyager_api.SchemeGroupVersion.WithResource(voyager_api.ResourceTypeIngress)},
		//{&voyager_api.Certificate{}, voyager_api.SchemeGroupVersion.WithResource(voyager_api.ResourceTypeCertificate)},
		//
		//// stash
		//{&stash_api.Restic{}, stash_api.SchemeGroupVersion.WithResource(stash_api.ResourceTypeRestic)},
		//{&stash_api.Recovery{}, stash_api.SchemeGroupVersion.WithResource(stash_api.ResourceTypeRecovery)},
		//
		//// searchlight
		//{&searchlight_api.ClusterAlert{}, searchlight_api.SchemeGroupVersion.WithResource(searchlight_api.ResourceTypeClusterAlert)},
		//{&searchlight_api.NodeAlert{}, searchlight_api.SchemeGroupVersion.WithResource(searchlight_api.ResourceTypeNodeAlert)},
		//{&searchlight_api.PodAlert{}, searchlight_api.SchemeGroupVersion.WithResource(searchlight_api.ResourceTypePodAlert)},
		//
		//// kubedb
		//{&kubedb_api.Postgres{}, kubedb_api.SchemeGroupVersion.WithResource(kubedb_api.ResourceTypePostgres)},
		//{&kubedb_api.Elasticsearch{}, kubedb_api.SchemeGroupVersion.WithResource(kubedb_api.ResourceTypeElasticsearch)},
		//{&kubedb_api.MySQL{}, kubedb_api.SchemeGroupVersion.WithResource(kubedb_api.ResourceTypeMySQL)},
		//{&kubedb_api.MongoDB{}, kubedb_api.SchemeGroupVersion.WithResource(kubedb_api.ResourceTypeMongoDB)},
		//{&kubedb_api.Redis{}, kubedb_api.SchemeGroupVersion.WithResource(kubedb_api.ResourceTypeRedis)},
		//{&kubedb_api.Memcached{}, kubedb_api.SchemeGroupVersion.WithResource(kubedb_api.ResourceTypeMemcached)},
		//{&kubedb_api.Snapshot{}, kubedb_api.SchemeGroupVersion.WithResource(kubedb_api.ResourceTypeSnapshot)},
		//{&kubedb_api.DormantDatabase{}, kubedb_api.SchemeGroupVersion.WithResource(kubedb_api.ResourceTypeDormantDatabase)},

		//// coreos-prometheus
		//{&prom.Prometheus{}, prom_util.SchemeGroupVersion.WithResource(prom.PrometheusesKind)},
		//{&prom.ServiceMonitor{}, prom_util.SchemeGroupVersion.WithResource(prom.ServiceMonitorsKind)},
		//{&prom.Alertmanager{}, prom_util.SchemeGroupVersion.WithResource(prom.AlertmanagersKind)},
	}

	for _, tt := range data {
		gvr, err := discovery.DetectResource(restmapper, tt.in)
		if err != nil {
			t.Error(err)
		}
		if gvr != tt.out {
			t.Errorf("Failed to DetectResource: expected %+v, got %+v", tt.out, gvr)
		}
	}
}
