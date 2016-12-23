package dns

import (
	"errors"
	"fmt"
	"os"
	"strings"

	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

const (
	PodNamespace string = "POD_NAMESPACE"
)

func GetServiceClusterIP(client clientset.Interface, prefix, hostname string) (string, error) {
	serviceName, namespace, err := splitHostname(hostname)
	if err != nil {
		return "", err
	}
	service, err := client.Core().Services(namespace).Get(serviceName)
	if err != nil {
		return "", err
	}

	sKey := fmt.Sprintf("%s_SERVICE_NAME", strings.ToUpper(prefix))
	os.Setenv(sKey, serviceName)
	nKey := fmt.Sprintf("%s_SERVICE_NAMESPACE", strings.ToUpper(prefix))
	os.Setenv(nKey, namespace)
	return service.Spec.ClusterIP, nil
}

func splitHostname(hostName string) (string, string, error) {
	parts := strings.Split(hostName, ".")
	if len(parts) == 1 {
		namespace := os.Getenv(PodNamespace)
		if namespace != "" {
			return parts[0], namespace, nil
		}
		return "", "", errors.New("Kubernetes namespace not found in ENV")
	} else if len(parts) == 2 {
		return parts[0], parts[1], nil
	}
	return "", "", fmt.Errorf(`Invalid hostname "%v"`, hostName)
}
