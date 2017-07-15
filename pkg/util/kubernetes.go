package util

import (
	"strconv"

	searchlight "github.com/appscode/searchlight/api"
	stash "github.com/appscode/stash/api"
	voyager "github.com/appscode/voyager/api"
	kubedb "github.com/k8sdb/apimachinery/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	batchv2alpha1 "k8s.io/client-go/pkg/apis/batch/v2alpha1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	rbacv1alpha1 "k8s.io/client-go/pkg/apis/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/client-go/pkg/apis/rbac/v1beta1"
	storagev1 "k8s.io/client-go/pkg/apis/storage/v1"
	storagev1beta1 "k8s.io/client-go/pkg/apis/storage/v1beta1"
)

func IsPreferredAPIResource(kubeClient clientset.Interface, groupVersion, kind string) bool {
	if resourceList, err := kubeClient.Discovery().ServerPreferredResources(); err == nil {
		for _, resources := range resourceList {
			if resources.GroupVersion != groupVersion {
				continue
			}
			for _, resource := range resources.APIResources {
				if resources.GroupVersion == groupVersion && resource.Kind == kind {
					return true
				}
			}
		}
	}
	return false
}

func GetBool(m map[string]string, key string) (bool, error) {
	if m == nil {
		return false, nil
	}
	return strconv.ParseBool(m[key])
}

func GetString(m map[string]string, key string) string {
	if m == nil {
		return ""
	}
	return m[key]
}

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	switch v.(type) {
	case *metav1.APIResourceList:
		return apiv1.SchemeGroupVersion.WithKind("APIResourceList")
	case *apiv1.Pod:
		return apiv1.SchemeGroupVersion.WithKind("Pod")
	case *apiv1.ReplicationController:
		return apiv1.SchemeGroupVersion.WithKind("ReplicationController")
	case *apiv1.ConfigMap:
		return apiv1.SchemeGroupVersion.WithKind("ConfigMap")
	case *apiv1.Secret:
		return apiv1.SchemeGroupVersion.WithKind("Secret")
	case *apiv1.Service:
		return apiv1.SchemeGroupVersion.WithKind("Service")
	case *apiv1.PersistentVolumeClaim:
		return apiv1.SchemeGroupVersion.WithKind("PersistentVolumeClaim")
	case *apiv1.PersistentVolume:
		return apiv1.SchemeGroupVersion.WithKind("PersistentVolume")
	case *apiv1.Node:
		return apiv1.SchemeGroupVersion.WithKind("Node")
	case *apiv1.ServiceAccount:
		return apiv1.SchemeGroupVersion.WithKind("ServiceAccount")
	case *apiv1.Namespace:
		return apiv1.SchemeGroupVersion.WithKind("Namespace")
	case *apiv1.Endpoints:
		return apiv1.SchemeGroupVersion.WithKind("Endpoints")
	case *apiv1.ComponentStatus:
		return apiv1.SchemeGroupVersion.WithKind("ComponentStatus")
	case *apiv1.LimitRange:
		return apiv1.SchemeGroupVersion.WithKind("LimitRange")
	case *apiv1.Event:
		return apiv1.SchemeGroupVersion.WithKind("Event")
	case *extensions.Ingress:
		return extensions.SchemeGroupVersion.WithKind("Ingress")
	case *extensions.DaemonSet:
		return extensions.SchemeGroupVersion.WithKind("DaemonSet")
	case *extensions.ReplicaSet:
		return extensions.SchemeGroupVersion.WithKind("ReplicaSet")
	case *extensions.Deployment:
		return extensions.SchemeGroupVersion.WithKind("Deployment")
	case *extensions.ThirdPartyResource:
		return extensions.SchemeGroupVersion.WithKind("ThirdPartyResource")
	case *apps.StatefulSet:
		return apps.SchemeGroupVersion.WithKind("StatefulSet")
	case *apps.Deployment:
		return apps.SchemeGroupVersion.WithKind("Deployment")
	case *storagev1beta1.StorageClass:
		return storagev1beta1.SchemeGroupVersion.WithKind("StorageClass")
	case *storagev1.StorageClass:
		return storagev1.SchemeGroupVersion.WithKind("StorageClass")
	case *rbacv1alpha1.Role:
		return rbacv1alpha1.SchemeGroupVersion.WithKind("Role")
	case *rbacv1alpha1.RoleBinding:
		return rbacv1alpha1.SchemeGroupVersion.WithKind("RoleBinding")
	case *rbacv1alpha1.ClusterRole:
		return rbacv1alpha1.SchemeGroupVersion.WithKind("ClusterRole")
	case *rbacv1alpha1.ClusterRoleBinding:
		return rbacv1alpha1.SchemeGroupVersion.WithKind("ClusterRoleBinding")
	case *rbacv1beta1.Role:
		return rbacv1beta1.SchemeGroupVersion.WithKind("Role")
	case *rbacv1beta1.RoleBinding:
		return rbacv1beta1.SchemeGroupVersion.WithKind("RoleBinding")
	case *rbacv1beta1.ClusterRole:
		return rbacv1beta1.SchemeGroupVersion.WithKind("ClusterRole")
	case *rbacv1beta1.ClusterRoleBinding:
		return rbacv1beta1.SchemeGroupVersion.WithKind("ClusterRoleBinding")
	case *batchv2alpha1.CronJob:
		return batchv2alpha1.SchemeGroupVersion.WithKind("CronJob")
	case *batchv1.Job:
		return batchv1.SchemeGroupVersion.WithKind("Job")
	case *searchlight.ClusterAlert:
		return searchlight.V1alpha1SchemeGroupVersion.WithKind("ClusterAlert")
	case *searchlight.NodeAlert:
		return searchlight.V1alpha1SchemeGroupVersion.WithKind("NodeAlert")
	case *searchlight.PodAlert:
		return searchlight.V1alpha1SchemeGroupVersion.WithKind("PodAlert")
	case *voyager.Ingress:
		return voyager.V1beta1SchemeGroupVersion.WithKind("Ingress")
	case *voyager.Certificate:
		return voyager.V1beta1SchemeGroupVersion.WithKind("Certificate")
	case *kubedb.Postgres:
		return kubedb.V1alpha1SchemeGroupVersion.WithKind("Postgres")
	case *kubedb.Elastic:
		return kubedb.V1alpha1SchemeGroupVersion.WithKind("Elastic")
	case *kubedb.Snapshot:
		return kubedb.V1alpha1SchemeGroupVersion.WithKind("Snapshot")
	case *kubedb.DormantDatabase:
		return kubedb.V1alpha1SchemeGroupVersion.WithKind("DormantDatabase")
	case *stash.Restic:
		return stash.V1alpha1SchemeGroupVersion.WithKind("Restic")
	default:
		return schema.GroupVersionKind{}
	}
}
