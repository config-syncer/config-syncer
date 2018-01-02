---
title: API Server
description: API Server
menu:
  product_kubed_0.3.1:
    identifier: tutorials-apiserver
    name: API Server
    parent: tutorials
    weight: 10
product_name: kubed
menu_name: product_kubed_0.3.1
section_menu_id: tutorials
---

> New to Kubed? Please start [here](/docs/guides/README.md).

# Kubed API Server
Kubed includes an api server. It has 2 categories of endpoints:
 - Search objects
 - Reverse Lookup

---

Kubed api server is under active development and expected to change in future. We are also [exploring](https://github.com/appscode/kubed/issues/19) the idea of turning this into a UAS.

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
| `apiServer.enableSearchIndex` | `Optional`. If set to `true`, builds a search index for Kubernetes api objects using [bleve](https://github.com/blevesearch/bleve). |

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

Now, deploy Kubed operator in your cluster following the steps [here](/docs/setup/install.md). Once the operator pod is running, go to the next section.


## Using Kubed API Server
In this section, we will show how you can use the kubed api server.


### Search Kubernetes objects
To search for Kubernetes objects, use the `/search` URL of Kubed api server. 

```console
$ kubectl get pods -n kube-system
NAME                              READY     STATUS    RESTARTS   AGE
kube-addon-manager-minikube       1/1       Running   0          33m
kube-dns-1301475494-hglm0         3/3       Running   0          33m
kubed-operator-3234987584-sbgrf   1/1       Running   0          19s
kubernetes-dashboard-l8vlj        1/1       Running   0          33m

$ kubectl port-forward $(kubectl get pods --all-namespaces -l app=kubed -o jsonpath={.items[0].metadata.name}) -n kube-system 8080
Forwarding from 127.0.0.1:8080 -> 8080
E0727 03:51:10.186041   22995 portforward.go:212] Unable to create listener: Error listen tcp6 [::1]:8080: bind: cannot assign requested address
Handling connection for 8080

# in a separate terminal window
$ curl http://127.0.0.1:8080/search?q=dashboard > ./docs/examples/apiserver/search-result.json
```

Now, open the URL [http://127.0.0.1:8080/search?q=dashboard](http://127.0.0.1:8080/search?q=dashboard) in your browser.


## Reverse Lookup of Objects
Sometimes you may want to know which [Prometheus stores metrics for a given Pod X](https://github.com/coreos/prometheus-operator/issues/230). Using reverse indices maintained by Kubed, answering questions like this become easier. Kubed maintains the following types of reverse indices:
 - List all Services for a given Pod.
 - List all ServiceMonitors for a given Service.
 - List all Prometheus objects for a given ServiceMonitor.

Since these indices are built using watchers, they always lag behind the current truth. But they work well for practical purposes.

```console
$ kubectl get pods -n kube-system
NAME                              READY     STATUS    RESTARTS   AGE
kube-addon-manager-minikube       1/1       Running   0          33m
kube-dns-1301475494-hglm0         3/3       Running   0          33m
kubed-operator-3234987584-sbgrf   1/1       Running   0          19s
kubernetes-dashboard-l8vlj        1/1       Running   0          33m

$ kubectl port-forward $(kubectl get pods --all-namespaces -l app=kubed -o jsonpath={.items[0].metadata.name}) -n kube-system 8080
Forwarding from 127.0.0.1:8080 -> 8080
E0727 03:51:10.186041   22995 portforward.go:212] Unable to create listener: Error listen tcp6 [::1]:8080: bind: cannot assign requested address
Handling connection for 8080

# in a separate terminal window
$ curl http://127.0.0.1:8080/api/v1/namespaces/kube-system/pods/kubernetes-dashboard-l8vlj/services > ./docs/examples/apiserver/pod-2-svc.json
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1283  100  1283    0     0  89128      0 --:--:-- --:--:-- --:--:-- 91642
```

Now, open the URL [http://127.0.0.1:8080/api/v1/namespaces/kube-system/pods/{pod-in-kube-system}/services](http://127.0.0.1:8080/api/v1/namespaces/kube-system/pods/{pod-in-kube-system}/services) in your browser.


## Supported Kubernetes Objects
Following Kubernetes objects are supported by search index:
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

To add support for additional object types, please [file an issue](https://github.com/appscode/kubed/issues/new?title=Support+Object+Kind+[xyz]+in+SearchEngine).


## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, uninstall Kubed operator following the steps [here](/docs/setup/uninstall.md).


## Next Steps
 - Learn how to use Kubed to take periodic snapshots of a Kubernetes cluster [here](/docs/guides/cluster-snapshot.md).
 - To setup a recycle bin for deleted and/or updated Kubernetes objects, please visit [here](/docs/guides/recycle-bin.md).
 - Need to keep some configuration synchronized across namespaces? Try [Kubed config syncer](/docs/guides/config-syncer.md).
 - Want to keep an eye on your cluster with automated notifications? Setup Kubed [event forwarder](/docs/guides/event-forwarder.md).
 - Out of disk space because of too much logs in Elasticsearch or metrics in InfluxDB? Configure [janitors](/docs/guides/janitors.md) to delete old data.
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/ROADMAP.md).
 - Want to hack on Kubed? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
