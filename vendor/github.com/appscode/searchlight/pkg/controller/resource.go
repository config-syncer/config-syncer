package controller

import (
	"fmt"
	"os"

	"github.com/appscode/errors"
	"github.com/appscode/k8s-addons/pkg/events"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/controller/types"
	kapi "k8s.io/kubernetes/pkg/api"
	k8error "k8s.io/kubernetes/pkg/api/errors"
)

func (b *IcingaController) IsObjectExists() error {
	log.Infoln("Checking Kubernetes Object existance", b.ctx.Resource.ObjectMeta)
	b.parseAlertOptions()

	var err error
	switch b.ctx.ObjectType {
	case events.Service.String():
		_, err = b.ctx.KubeClient.Core().Services(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName)
	case events.RC.String():
		_, err = b.ctx.KubeClient.Core().ReplicationControllers(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName)
	case events.DaemonSet.String():
		_, err = b.ctx.KubeClient.Extensions().DaemonSets(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName)
	case events.Deployments.String():
		_, err = b.ctx.KubeClient.Extensions().Deployments(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName)
	case events.StatefulSet.String():
		_, err = b.ctx.KubeClient.Apps().StatefulSets(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName)
	case events.ReplicaSet.String():
		_, err = b.ctx.KubeClient.Extensions().ReplicaSets(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName)
	case events.Pod.String():
		_, err = b.ctx.KubeClient.Core().Pods(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName)
	case events.Node.String():
		if b.ctx.ObjectName == "" {
			return nil
		}
		if _, err = b.ctx.KubeClient.Core().Nodes().Get(b.ctx.ObjectName); err != nil {
			if k8error.IsNotFound(err) {
				return errors.New(fmt.Sprintf(`Node "%s" not found`, b.ctx.ObjectName)).NotFound()
			}
			return errors.New().WithCause(err)
		}
	case events.Cluster.String():
		return nil
	default:
		return errors.New(fmt.Sprintf(`Invalid Object Type "%s"`, b.ctx.ObjectType)).InvalidData()
	}

	if err != nil {
		if k8error.IsNotFound(err) {
			return errors.New(fmt.Sprintf(`Kubernetes Object "%s" of kind "%s" in namespace "%s" not found`, b.ctx.ObjectName, b.ctx.ObjectType, b.ctx.Resource.Namespace)).NotFound()
		}
		return errors.New().WithCause(err)
	}

	return nil
}

func (b *IcingaController) getParentsForPod(o interface{}) []*types.Ancestors {
	pod := o.(*kapi.Pod)
	result := make([]*types.Ancestors, 0)

	svc, err := b.ctx.Storage.ServiceStore.GetPodServices(pod)
	if err == nil {
		names := make([]string, 0)
		for _, s := range svc {
			names = append(names, s.Name)
		}
		result = append(result, &types.Ancestors{
			Type:  events.Service.String(),
			Names: names,
		})
	}

	rc, err := b.ctx.Storage.RcStore.GetPodControllers(pod)
	if err == nil {
		names := make([]string, 0)
		for _, s := range rc {
			names = append(names, s.Name)
		}
		result = append(result, &types.Ancestors{
			Type:  events.RC.String(),
			Names: names,
		})
	}

	rs, err := b.ctx.Storage.ReplicaSetStore.GetPodReplicaSets(pod)
	if err == nil {
		names := make([]string, 0)
		for _, s := range rs {
			names = append(names, s.Name)
		}
		result = append(result, &types.Ancestors{
			Type:  events.ReplicaSet.String(),
			Names: names,
		})
	}

	ps, err := b.ctx.Storage.StatefulSetStore.GetPodStatefulSets(pod)
	if err == nil {
		names := make([]string, 0)
		for _, s := range ps {
			names = append(names, s.Name)
		}
		result = append(result, &types.Ancestors{
			Type:  events.StatefulSet.String(),
			Names: names,
		})
	}

	ds, err := b.ctx.Storage.DaemonSetStore.GetPodDaemonSets(pod)
	if err == nil {
		names := make([]string, 0)
		for _, s := range ds {
			names = append(names, s.Name)
		}
		result = append(result, &types.Ancestors{
			Type:  events.DaemonSet.String(),
			Names: names,
		})
	}
	return result
}

func (b *IcingaController) checkIcingaAvailability() bool {
	log.Debugln("Checking Icinga client")
	if b.ctx.IcingaClient == nil {
		return false
	}
	resp := b.ctx.IcingaClient.Check().Get([]string{}).Do()
	if resp.Status != 200 {
		return false
	}
	return true
}

func (b *IcingaController) checkPodIPAvailability(podName, namespace string) (bool, error) {
	log.Debugln("Checking pod IP")
	pod, err := b.ctx.KubeClient.Core().Pods(namespace).Get(podName)
	if err != nil {
		return false, errors.New().WithCause(err).Internal()
	}
	if pod.Status.PodIP == "" {
		return false, nil
	}
	return true, nil
}

func checkIcingaService(serviceName, namespace string) bool {
	icingaService := os.Getenv("ICINGA_SERVICE_NAME")
	if serviceName != icingaService {
		return false
	}
	icingaServiceNamespace := os.Getenv("ICINGA_SERVICE_NAMESPACE")
	if namespace != icingaServiceNamespace {
		return false
	}
	return true
}
