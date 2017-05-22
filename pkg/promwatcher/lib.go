package watcher

import (
	"errors"
	"fmt"
	"time"

	aci "github.com/appscode/k8s-addons/api"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	kapi "k8s.io/kubernetes/pkg/api"
	k8serr "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/resource"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/util/intstr"
)

func (p *PromWatcher) createGoverningService(prometheus *pcm.Prometheus) (*kapi.Service, error) {
	service, err := p.Client.Core().Services(prometheus.Namespace).Get(keyPrometheus)
	if err != nil {
		if !k8serr.IsNotFound(err) {
			return nil, err
		}
	} else {
		return service, nil
	}

	labels := prometheus.Labels
	labels["app"] = keyPrometheus

	service = &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Name:      keyPrometheus,
			Namespace: prometheus.Namespace,
		},
		Spec: kapi.ServiceSpec{
			Selector: labels,
			Ports: []kapi.ServicePort{
				{
					Name:       "web",
					Protocol:   kapi.ProtocolTCP,
					Port:       9090,
					TargetPort: intstr.FromInt(9090),
				},
			},
			Type:      kapi.ServiceTypeClusterIP,
			ClusterIP: kapi.ClusterIPNone,
		},
	}

	return p.Client.Core().Services(service.Namespace).Create(service)
}

const authProxyImage string = "appscode/authui-proxy:0.3.0"

func (p *PromWatcher) createProxyDeployment(prometheusName, backendServiceName, namespace string) (*extensions.Deployment, error) {
	deploymentName := fmt.Sprintf("%v-%v-proxy", keyPrometheus, prometheusName)
	labels := map[string]string{
		"proxy.prometheus/name": deploymentName,
	}

	resourceQuanty, err := resource.ParseQuantity("100m")
	if err != nil {
		return nil, err
	}

	petAddress := fmt.Sprintf("%v-%v-0.%v.%v.svc.cluster.local",
		keyPrometheus, prometheusName, keyPrometheus, namespace)

	deployment := &extensions.Deployment{
		ObjectMeta: kapi.ObjectMeta{
			Name:   deploymentName,
			Labels: labels,
		},
		Spec: extensions.DeploymentSpec{
			Template: kapi.PodTemplateSpec{
				ObjectMeta: kapi.ObjectMeta{
					Labels: labels,
				},
				Spec: kapi.PodSpec{
					Containers: []kapi.Container{
						{
							Name:  "prometheus-proxy",
							Image: authProxyImage,
							Args: []string{
								fmt.Sprintf("--app-url=http://%v:9090", petAddress),
								"--home-path=/",
								"--base-prefix=/",
								"--proxy-port=9090",
								"--v=3",
							},
							Resources: kapi.ResourceRequirements{
								Limits: kapi.ResourceList{
									kapi.ResourceCPU: resourceQuanty,
								},
								Requests: kapi.ResourceList{
									kapi.ResourceCPU: resourceQuanty,
								},
							},
							Ports: []kapi.ContainerPort{
								{
									Name:          "ui",
									Protocol:      kapi.ProtocolTCP,
									ContainerPort: 9090,
								},
							},
							VolumeMounts: []kapi.VolumeMount{
								{
									Name:      "appscode-cluster-metadata",
									MountPath: "/var/run/config/appscode",
								},
							},
						},
					},
					Volumes: []kapi.Volume{
						{
							Name: "appscode-cluster-metadata",
							VolumeSource: kapi.VolumeSource{
								ConfigMap: &kapi.ConfigMapVolumeSource{
									LocalObjectReference: kapi.LocalObjectReference{
										Name: "cluster-metadata",
									},
								},
							},
						},
					},
				},
			},
			Replicas: 1,
		},
	}

	return p.Client.Extensions().Deployments(namespace).Create(deployment)
}

func (p *PromWatcher) deleteProxyDeployment(prometheusName, namespace string) error {
	deploymentName := fmt.Sprintf("%v-%v-proxy", keyPrometheus, prometheusName)
	deployment, err := p.Client.Extensions().Deployments(namespace).Get(deploymentName)
	if err != nil {
		return err
	}
	deployment.Spec.Replicas = 0
	if _, err := p.Client.Extensions().Deployments(deployment.Namespace).Update(deployment); err != nil {
		return err
	}

	labelSelector := labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels)
	time.Sleep(time.Second * 30)
	check := 0
	for {
		podList, err := p.Client.Core().Pods(deployment.Namespace).List(kapi.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			return err
		}
		if len(podList.Items) == 0 {
			break
		}

		if check > 6 {
			return errors.New("Fail to delete Deployment Pods")
		}
		time.Sleep(time.Second * 10)
		check++
	}

	// Delete Deployment
	return p.Client.Extensions().Deployments(deployment.Namespace).Delete(deployment.Name, nil)
}

func (p *PromWatcher) createProxyService(prometheusName string, deployment *extensions.Deployment) error {
	labels := deployment.Labels

	serviceName := fmt.Sprintf("%v-%v-proxy", keyPrometheus, prometheusName)
	service := &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Name:      serviceName,
			Namespace: deployment.Namespace,
		},
		Spec: kapi.ServiceSpec{
			Selector: labels,
			Ports: []kapi.ServicePort{
				{
					Name:       "ui",
					Protocol:   kapi.ProtocolTCP,
					Port:       9090,
					TargetPort: intstr.FromInt(9090),
				},
			},
			Type: kapi.ServiceTypeNodePort,
		},
	}
	_, err := p.Client.Core().Services(service.Namespace).Create(service)
	return err
}

func (p *PromWatcher) deleteProxyService(prometheusName, namespace string) error {
	proxyService := fmt.Sprintf("%v-%v-proxy", keyPrometheus, prometheusName)
	return p.Client.Core().Services(namespace).Delete(proxyService, nil)
}

const (
	ingressName      = "default-lb"
	ingressNamespace = "appscode"
)

func (p *PromWatcher) createIngressRule(prometheusName, namespace string) error {
	ingress, err := p.AppsCodeExtensionClient.Ingress(ingressNamespace).Get(ingressName)
	if err != nil {
		return err
	}

	prefix := fmt.Sprintf("/%v-%v.%v", keyPrometheus, prometheusName, namespace)
	proxyServiceName := fmt.Sprintf("%v-%v-proxy", keyPrometheus, prometheusName)

	rule := aci.ExtendedIngressRule{
		ExtendedIngressRuleValue: aci.ExtendedIngressRuleValue{
			HTTP: &aci.HTTPExtendedIngressRuleValue{
				Paths: []aci.HTTPExtendedIngressPath{
					{
						Path: prefix,
						Backend: aci.ExtendedIngressBackend{
							ServiceName: fmt.Sprintf("%v.%v", proxyServiceName, namespace),
							ServicePort: intstr.FromInt(9090),
						},
					},
				},
			},
		},
	}

	ingress.Spec.Rules = append(ingress.Spec.Rules, rule)

	if _, err := p.AppsCodeExtensionClient.Ingress(ingress.Namespace).Update(ingress); err != nil {
		return err
	}
	return nil
}

func (p *PromWatcher) deleteIngressRule(prometheusName, namespace string) error {
	ingress, err := p.AppsCodeExtensionClient.Ingress(ingressNamespace).Get(ingressName)
	if err != nil {
		return err
	}

	prefix := fmt.Sprintf("/%v-%v.%v", keyPrometheus, prometheusName, namespace)

	for i, rule := range ingress.Spec.Rules {
		found := false
		for _, path := range rule.HTTP.Paths {
			if path.Path == prefix {
				found = true
				break
			}
		}
		if found {
			ingress.Spec.Rules = append(ingress.Spec.Rules[:i], ingress.Spec.Rules[i+1:]...)
			break
		}
	}
	return nil
}
