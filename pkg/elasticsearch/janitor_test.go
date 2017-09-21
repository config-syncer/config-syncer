package es

import (
	"testing"
	"time"

	"github.com/appscode/kubed/pkg/config"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestEsJanitor(t *testing.T) {
	c, err := clientcmd.BuildConfigFromFlags("", "")
	assert.Nil(t, err)

	kubeClient := clientset.NewForConfigOrDie(c)

	var authInfo *config.JanitorAuthInfo
	secret, err := kubeClient.CoreV1().Secrets("").Get("", metav1.GetOptions{})
	assert.Nil(t, err)

	authInfo = config.LoadJanitorAuthInfo(secret.Data)

	esSpec := config.ElasticsearchSpec{
		Endpoint:       "https://localhost:32317",
		LogIndexPrefix: "logstash-",
	}

	janitor := Janitor{Spec: esSpec, AuthInfo: authInfo, TTL: time.Minute * 10}
	err = janitor.Cleanup()
	assert.Nil(t, err)
}
