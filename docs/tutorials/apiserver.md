> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Kubed API Server
Kubed provides a recycle bin for deleted and/or updated Kubernetes objects. Once activated, any deleted and/or updated object is stored in YAML format in folder mounted inside Kubed pod. This tutorial will show you how to use Kubed to setup a recycle bin for Kubernetes cluster objects.

---

Kubed api server is under active development and expected to change in future.

---

## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).


## Deploy Kubed
To enable config syncer, you need a cluster config like below.

```yaml
$ cat ./docs/examples/apiserver/config.yaml

apiServer:
  address: :8080
  enableReverseIndex: true
  enableSearchIndex: true
```

| Key                            | Description                                                                                    |
|--------------------------------|------------------------------------------------------------------------------------------------|
| `apiServer.address`            | `Optional`. Address of the Kubed API Server (can be overridden by `kubed run --address` flag). |
| `apiServer.enableReverseIndex` | `Optional`. If set to `true`, builds a reverse index                                           |
| `apiServer.enableReverseIndex` | `Optional`. If set to `true`, builds a search index for Kubernetes api objects using [bleve](https://github.com/blevesearch/bleve). |

Now, create a Secret with the Kubed cluster config under `config.yaml` key.

```yaml
$ kubectl create secret generic kubed-config -n kube-system \
    --from-file=./docs/examples/apiserver/config.yaml
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

$ kubectl get secret kubed-config -n kube-system -o yaml
apiVersion: v1
data:
  config.yaml: YXBpU2VydmVyOgogIGFkZHJlc3M6IDo4MDgwCiAgZW5hYmxlUmV2ZXJzZUluZGV4OiB0cnVlCiAgZW5hYmxlU2VhcmNoSW5kZXg6IHRydWUK
kind: Secret
metadata:
  creationTimestamp: 2017-07-27T10:47:41Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "2187"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: 039ab470-72b9-11e7-a1f7-080027df84b0
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/install.md). Once the operator pod is running, go to the next section.


## Using Kubed API Server
















```console
$ kubectl get pods -n kube-system
NAME                              READY     STATUS    RESTARTS   AGE
kube-addon-manager-minikube       1/1       Running   0          33m
kube-dns-1301475494-hglm0         3/3       Running   0          33m
kubed-operator-3234987584-sbgrf   1/1       Running   0          19s
kubernetes-dashboard-l8vlj        1/1       Running   0          33m


$ kubectl port-forward $(kubectl get pods --all-namespaces -l app=kubed -o jsonpath={.items[0].metadata.name}) -n kube-system 56790
Forwarding from 127.0.0.1:56790 -> 56790
E0727 03:50:34.668103   22871 portforward.go:212] Unable to create listener: Error listen tcp6 [::1]:56790: bind: cannot assign requested address
Handling connection for 56790
^C⏎


$ kubectl port-forward $(kubectl get pods --all-namespaces -l app=kubed -o jsonpath={.items[0].metadata.name}) -n kube-system 8080
Forwarding from 127.0.0.1:8080 -> 8080
E0727 03:51:10.186041   22995 portforward.go:212] Unable to create listener: Error listen tcp6 [::1]:8080: bind: cannot assign requested address
Handling connection for 8080
^C⏎

$ curl http://127.0.0.1:8080/search?q=dashboard > ./docs/examples/apiserver/search-result.json


                                                                                                                                                             
$ curl http://127.0.0.1:8080/api/v1/namespaces/kube-system/pods/kubernetes-dashboard-l8vlj/services > ./docs/examples/apiserver/pod-2-svc.json
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1283  100  1283    0     0  89128      0 --:--:-- --:--:-- --:--:-- 91642
```


		router.Get("/api/v1/namespaces/:namespace/:resource/:name/services", http.HandlerFunc(op.ReverseIndex.Service.ServeHTTP))
		if util.IsPreferredAPIResource(op.KubeClient, prom.TPRGroup+"/"+prom.TPRVersion, prom.TPRServiceMonitorsKind) {
			// Add Indexer only if Server support this resource
			router.Get("/apis/"+prom.TPRGroup+"/"+prom.TPRVersion+"/namespaces/:namespace/:resource/:name/"+prom.TPRServiceMonitorName, http.HandlerFunc(op.ReverseIndex.ServiceMonitor.ServeHTTP))
		}
		if util.IsPreferredAPIResource(op.KubeClient, prom.TPRGroup+"/"+prom.TPRVersion, prom.TPRPrometheusesKind) {
			// Add Indexer only if Server support this resource
			router.Get("/apis/"+prom.TPRGroup+"/"+prom.TPRVersion+"/namespaces/:namespace/:resource/:name/"+prom.TPRPrometheusName, http.HandlerFunc(op.ReverseIndex.Prometheus.ServeHTTP))
		}


https://github.com/coreos/prometheus-operator/issues/230




## Kubed Metrics Server
kubed exposes Prometheus ready metrics via an endpoint running 

```console
$ kubectl get pods -n kube-system
NAME                              READY     STATUS    RESTARTS   AGE
kube-addon-manager-minikube       1/1       Running   0          33m
kube-dns-1301475494-hglm0         3/3       Running   0          33m
kubed-operator-3234987584-sbgrf   1/1       Running   0          19s
kubernetes-dashboard-l8vlj        1/1       Running   0          33m


$ kubectl port-forward $(kubectl get pods --all-namespaces -l app=kubed -o jsonpath={.items[0].metadata.name}) -n kube-system 56790
Forwarding from 127.0.0.1:56790 -> 56790
E0727 03:50:34.668103   22871 portforward.go:212] Unable to create listener: Error listen tcp6 [::1]:56790: bind: cannot assign requested address
Handling connection for 56790
^C⏎


$ kubectl port-forward $(kubectl get pods --all-namespaces -l app=kubed -o jsonpath={.items[0].metadata.name}) -n kube-system 8080
Forwarding from 127.0.0.1:8080 -> 8080
E0727 03:51:10.186041   22995 portforward.go:212] Unable to create listener: Error listen tcp6 [::1]:8080: bind: cannot assign requested address
Handling connection for 8080
^C⏎

$ curl http://127.0.0.1:8080/search?q=dashboard > ./docs/examples/apiserver/search-result.json


                                                                                                                                                             
$ curl http://127.0.0.1:8080/api/v1/namespaces/kube-system/pods/kubernetes-dashboard-l8vlj/services > ./docs/examples/apiserver/pod-2-svc.json
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1283  100  1283    0     0  89128      0 --:--:-- --:--:-- --:--:-- 91642
```





## Using Kubed API Server
In this tutorial, a ConfigMap will be used to show how recycle bin feature can be used.

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create namespace demo
namespace "demo" created

~ $ kubectl get namespaces
NAME          STATUS    AGE
default       Active    6h
kube-public   Active    6h
kube-system   Active    6h
demo          Active    4m
```

Create a ConfigMap called `omni` in the `demo` namespace.

```console
$ kubectl create configmap omni -n demo --from-literal=hello=world
configmap "omni" created
```
```yaml
$ kubectl get configmaps omni -n demo -o yaml
apiVersion: v1
data:
  hello: world
kind: ConfigMap
metadata:
  creationTimestamp: 2017-07-26T19:18:40Z
  name: omni
  namespace: demo
  resourceVersion: "34414"
  selfLink: /api/v1/namespaces/demo/configmaps/omni
  uid: 3b77f592-7237-11e7-af79-08002738e55e
```

Now, delete the ConfigMap `omni`. Kubed operator pod will notice this and stored the deleted object in YAML format in a file matching the `selfLink` for that object inside the `recycleBin.path` folder.

```console
# Exec into kubed operator pod
$ kubectl exec -it $(kubectl get pods --all-namespaces -l app=kubed -o jsonpath={.items[0].metadata.name}) -n kube-system sh

# running inside kubed operator pod
/ # find /tmp/kubed/trash/
/tmp/kubed/trash/
/tmp/kubed/trash/api
/tmp/kubed/trash/api/v1
/tmp/kubed/trash/api/v1/namespaces
/tmp/kubed/trash/api/v1/namespaces/demo
/tmp/kubed/trash/api/v1/namespaces/demo/configmaps
/tmp/kubed/trash/api/v1/namespaces/demo/configmaps/omni.20170726T193302.yaml

/ # cat /tmp/kubed/trash/api/v1/namespaces/demo/configmaps/omni.20170726T193302.yaml
apiVersion: v1
data:
  hello: world
kind: ConfigMap
metadata:
  creationTimestamp: 2017-07-26T19:33:02Z
  name: omni
  namespace: demo
  resourceVersion: "35481"
  selfLink: /api/v1/namespaces/demo/configmaps/omni
  uid: 3d50fba0-7239-11e7-af79-08002738e55e
```

## Supported Kubernetes Objects
Following Kubernetes objects are supported by recycle bin:
- __v1:__
  - ComponentStatus
  - ConfigMap
  - Endpoints
  - Event
  - LimitRange
  - Namespace
  - Node
  - PersistentVolume
  - PersistentVolumeClaim
  - Pod
  - ReplicationController
  - Secret
  - Service
  - ServiceAccount
- __apps/v1beta1:__
  - Deployment
  - StatefulSet
- __batch/v1:__
  - Job
- __batch/v2alpha1:__
  - CronJob
- __extensions/v1beta1:__
  - DaemonSet
  - Deployment
  - Ingress
  - ReplicaSet
  - ThirdPartyResource
- __kubedb/v1alpha1:__
  - DormantDatabase
  - Elasticsearch
  - Postgres
  - Snapshot
- __monitoring.coreos.com:__
  - Prometheus
  - ServiceMonitor
- __rbac/v1alpha1:__
  - ClusterRole
  - ClusterRoleBinding
  - Role
  - RoleBinding
- __rbac/v1beta1:__
  - ClusterRole
  - ClusterRoleBinding
  - Role
  - RoleBinding
- __monitoring.appscode.com/v1alpha1:__
  - ClusterAlert
  - NodeAlert
  - PodAlert
- __stash.appscode.com/v1alpha1:__
  - Restic
- __storage/v1:__
  - StorageClass
- __storage/v1beta1:__
  - StorageClass
- __voyager.appscode.com/v1beta1:__
  - Certificate
  - Ingress

To add support for additional object types, please [file an issue](https://github.com/appscode/kubed/issues/new?title=Support+Object+Kind+[xyz]+in+RecycleBin). We are exploring ways to watch for any object deletion [here](https://github.com/appscode/kubed/issues/41).


## Using Persistent Storage
The installation scripts for Kubed mounts an `emptyDir` under `/tmp` path. This tutorial used `/tmp/kubed/trash` to store objects in recycle bin. If you want to use a [PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) to recycle bin data, mount a PV and updated the `recycleBin.path` accordingly.


## Disable Recycle Bin
If you would like to disable this feature, remove the `recyclebin` portion of your Kubed cluster config. Then update the `kubed-config` Secret and restart Kubed operator pod(s).


## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run the following commands:
```console
$ kubectl delete ns demo
namespace "demo" deleted
```

To uninstall Kubed operator, please follow the steps [here](/docs/uninstall.md).

## Next Steps
