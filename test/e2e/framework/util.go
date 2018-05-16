package framework

import (
	"path/filepath"

	"github.com/appscode/go/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	KubedTestConfigFileDir = filepath.Join(runtime.GOPath(), "src", "github.com", "appscode", "kubed", "test", "e2e", "config.yaml")
)

func deleteInBackground() *metav1.DeleteOptions {
	policy := metav1.DeletePropagationBackground
	return &metav1.DeleteOptions{PropagationPolicy: &policy}
}

func deleteInForeground() *metav1.DeleteOptions {
	policy := metav1.DeletePropagationForeground
	return &metav1.DeleteOptions{PropagationPolicy: &policy}
}
