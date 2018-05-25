package syncer

import (
	"strings"

	"github.com/appscode/go/types"
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kutil/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

type SyncOptions struct {
	NamespaceSelector *string // if nil, delete from cluster
	Contexts          sets.String
}

func GetSyncOptions(annotations map[string]string) SyncOptions {
	opts := SyncOptions{}
	if v, err := meta.GetStringValue(annotations, api.ConfigSyncKey); err == nil {
		if v == "true" {
			opts.NamespaceSelector = types.StringP(labels.Everything().String())
		} else {
			opts.NamespaceSelector = &v
		}
	}
	if contexts, _ := meta.GetStringValue(annotations, api.ConfigSyncContexts); contexts != "" {
		opts.Contexts = sets.NewString(strings.Split(contexts, ",")...)
	}
	return opts
}

func NamespacesForSelector(kubeClient kubernetes.Interface, selector string) (sets.String, error) {
	namespaces, err := kubeClient.CoreV1().Namespaces().List(metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}
	ns := sets.NewString()
	for _, obj := range namespaces.Items {
		ns.Insert(obj.Name)
	}
	return ns, nil
}
