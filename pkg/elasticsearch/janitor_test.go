package es

import (
	"testing"
	"time"

	apis "github.com/appscode/kubed/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestEsJanitor(t *testing.T) {
	c, err := clientcmd.BuildConfigFromFlags("", "")
	assert.Nil(t, err)

	kubeClient := kubernetes.NewForConfigOrDie(c)

	var authInfo *apis.JanitorAuthInfo
	secret, err := kubeClient.CoreV1().Secrets("").Get("", metav1.GetOptions{})
	assert.Nil(t, err)

	authInfo = apis.LoadJanitorAuthInfo(secret.Data)

	esSpec := apis.ElasticsearchSpec{
		Endpoint:       "https://localhost:32317",
		LogIndexPrefix: "logstash-",
	}

	janitor := Janitor{Spec: esSpec, AuthInfo: authInfo, TTL: time.Minute * 10}
	err = janitor.Cleanup()
	assert.Nil(t, err)
}
