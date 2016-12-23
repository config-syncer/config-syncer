package host

import (
	"github.com/appscode/errors"
	aci "github.com/appscode/k8s-addons/api"
	acs "github.com/appscode/k8s-addons/client/clientset"
	kapi "k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/selection"
	"k8s.io/kubernetes/pkg/util/sets"
)

const (
	TypeServices               = "services"
	TypeReplicationcontrollers = "replicationcontrollers"
	TypeDaemonsets             = "daemonsets"
	TypeStatefulSet            = "statefulsets"
	TypeReplicasets            = "replicasets"
	TypeDeployments            = "deployments"
	TypePods                   = "pods"
	TypeNodes                  = "nodes"
	TypeCluster                = "cluster"
)

func getLabels(client clientset.Interface, namespace, objectType, objectName string) (labels.Selector, error) {
	label := labels.NewSelector()
	labelsMap := make(map[string]string, 0)
	if objectType == TypeServices {
		service, err := client.Core().Services(namespace).Get(objectName)
		if err != nil {
			return nil, errors.New().WithCause(err).Internal()
		}
		labelsMap = service.Spec.Selector

	} else if objectType == TypeReplicationcontrollers {
		rc, err := client.Core().ReplicationControllers(namespace).Get(objectName)
		if err != nil {
			return nil, errors.New().WithCause(err).Internal()
		}
		labelsMap = rc.Spec.Selector
	} else if objectType == TypeDaemonsets {
		daemonSet, err := client.Extensions().DaemonSets(namespace).Get(objectName)
		if err != nil {
			return nil, errors.New().WithCause(err).Internal()
		}
		labelsMap = daemonSet.Spec.Selector.MatchLabels
	} else if objectType == TypeReplicasets {
		replicaSet, err := client.Extensions().ReplicaSets(namespace).Get(objectName)
		if err != nil {
			return nil, errors.New().WithCause(err).Internal()
		}
		labelsMap = replicaSet.Spec.Selector.MatchLabels
	} else if objectType == TypeStatefulSet {
		petSet, err := client.Apps().StatefulSets(namespace).Get(objectName)
		if err != nil {
			return nil, errors.New().WithCause(err).Internal()
		}
		labelsMap = petSet.Spec.Selector.MatchLabels
	} else if objectType == TypeDeployments {
		deployment, err := client.Extensions().Deployments(namespace).Get(objectName)
		if err != nil {
			return nil, errors.New().WithCause(err).Internal()
		}
		labelsMap = deployment.Spec.Selector.MatchLabels
	} else {
		return label, errors.New().WithMessage("Invalid kubernetes object type").BadRequest()
	}

	for key, value := range labelsMap {
		s := sets.NewString(value)
		ls, err := labels.NewRequirement(key, selection.Equals, s.List())
		if err != nil {
			return nil, errors.New().WithCause(err).Internal()
		}
		label = label.Add(*ls)
	}

	return label, nil
}

func GetPodList(client clientset.Interface, namespace, objectType, objectName string) ([]*KubeObjectInfo, error) {
	var podList []*KubeObjectInfo

	label, err := getLabels(client, namespace, objectType, objectName)
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}

	pods, err := client.Core().Pods(namespace).List(kapi.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}

	for _, pod := range pods.Items {
		podList = append(podList, &KubeObjectInfo{Name: pod.Name + "@" + namespace, IP: pod.Status.PodIP, GroupName: objectName, GroupType: objectType})
	}

	return podList, nil
}

func GetPod(client clientset.Interface, namespace, objectType, objectName, podName string) ([]*KubeObjectInfo, error) {
	var podList []*KubeObjectInfo
	pod, err := client.Core().Pods(namespace).Get(podName)
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}
	podList = append(podList, &KubeObjectInfo{Name: pod.Name + "@" + namespace, IP: pod.Status.PodIP, GroupName: objectName, GroupType: objectType})
	return podList, nil
}

func GetNodeList(client clientset.Interface, alertNamespace string) ([]*KubeObjectInfo, error) {
	var nodeList []*KubeObjectInfo
	nodes, err := client.Core().Nodes().List(kapi.ListOptions{LabelSelector: labels.Everything()})
	if err != nil {
		return nodeList, errors.New().WithCause(err).Internal()
	}
	for _, node := range nodes.Items {
		nodeIP := "127.0.0.1"
		for _, ip := range node.Status.Addresses {
			if ip.Type == internalIP {
				nodeIP = ip.Address
				break
			}
		}
		nodeList = append(nodeList, &KubeObjectInfo{Name: node.Name + "@" + alertNamespace, IP: nodeIP, GroupName: TypeNodes, GroupType: ""})
	}
	return nodeList, nil
}

func GetNode(client clientset.Interface, nodeName, alertNamespace string) ([]*KubeObjectInfo, error) {
	var nodeList []*KubeObjectInfo
	node := &kapi.Node{}
	node, err := client.Core().Nodes().Get(nodeName)
	if err != nil {
		return nodeList, errors.New().WithCause(err).Internal()
	}
	nodeIP := "127.0.0.1"
	for _, ip := range node.Status.Addresses {
		if ip.Type == internalIP {
			nodeIP = ip.Address
			break
		}
	}
	nodeList = append(nodeList, &KubeObjectInfo{Name: node.Name + "@" + alertNamespace, IP: nodeIP, GroupName: TypeNodes, GroupType: ""})
	return nodeList, nil
}

func GetAlertList(acExtClient acs.AppsCodeExtensionInterface, kubeClient clientset.Interface, namespace string, ls labels.Selector) ([]aci.Alert, error) {
	alerts := make([]aci.Alert, 0)
	if namespace != "" {
		alertList, err := acExtClient.Alert(namespace).List(kapi.ListOptions{LabelSelector: ls})
		if err != nil {
			return nil, errors.New().WithCause(err).Internal()
		}
		if len(alertList.Items) > 0 {
			alerts = append(alerts, alertList.Items...)
		}
	} else {
		namespaces, _ := kubeClient.Core().Namespaces().List(kapi.ListOptions{LabelSelector: labels.Everything()})
		for _, ns := range namespaces.Items {
			alertList, err := acExtClient.Alert(ns.Name).List(kapi.ListOptions{LabelSelector: ls})
			if err != nil {
				return nil, errors.New().WithCause(err).Internal()
			}
			if len(alertList.Items) > 0 {
				alerts = append(alerts, alertList.Items...)
			}
		}
	}

	return alerts, nil
}

func GetAlert(acExtClient acs.AppsCodeExtensionInterface, namespace, name string) (*aci.Alert, error) {
	return acExtClient.Alert(namespace).Get(name)
}

const (
	ObjectType = "alert.appscode.com/objectType"
	ObjectName = "alert.appscode.com/objectName"
	AppName    = "k8s-app"
	AlertApp   = "appscode-alert"
)

func GetLabelSelector(objectType, objectName string) (labels.Selector, error) {
	lb := labels.NewSelector()
	if objectType != "" {
		lsot, err := labels.NewRequirement(ObjectType, selection.Equals, sets.NewString(objectType).List())
		if err != nil {
			return lb, errors.New().WithCause(err).Internal()
		}
		lb = lb.Add(*lsot)
	}

	if objectName != "" {
		lson, err := labels.NewRequirement(ObjectName, selection.Equals, sets.NewString(objectName).List())
		if err != nil {
			return lb, errors.New().WithCause(err).Internal()
		}
		lb = lb.Add(*lson)
	}

	return lb, nil
}

type labelMap map[string]string

func (s labelMap) ObjectType() string {
	v, _ := s[ObjectType]
	return v
}

func (s labelMap) ObjectName() string {
	v, _ := s[ObjectName]
	return v
}

func (s labelMap) AppName() string {
	v, _ := s[AppName]
	return v
}

func GetObjectInfo(label map[string]string) (objectType string, objectName string) {
	opts := labelMap(label)
	objectType = opts.ObjectType()
	objectName = opts.ObjectName()
	return
}

func CheckAlertConfig(oldConfig, newConfig *aci.Alert) error {
	oldOpts := labelMap(oldConfig.ObjectMeta.Labels)
	newOpts := labelMap(newConfig.ObjectMeta.Labels)

	if newOpts.ObjectType() != oldOpts.ObjectType() {
		return errors.New("Kubernetes ObjectType mismatch")
	}

	if newOpts.ObjectName() != oldOpts.ObjectName() {
		return errors.New("Kubernetes ObjectName mismatch")
	}

	if newConfig.Spec.CheckCommand != oldConfig.Spec.CheckCommand {
		return errors.New("CheckCommand mismatch")
	}

	return nil
}

func IsIcingaApp(labels map[string]string) bool {
	opts := labelMap(labels)
	return opts.AppName() == AlertApp
}
