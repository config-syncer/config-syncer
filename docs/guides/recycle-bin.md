---
title: Recycle Bin
description: Recycle Bin
menu:
  product_kubed_0.3.1:
    identifier: tutorials-recycle-bin
    name: Recycle Bin
    parent: tutorials
    weight: 35
product_name: kubed
menu_name: product_kubed_0.3.1
section_menu_id: tutorials
---

> New to Kubed? Please start [here](/docs/guides/README.md).

# Kubernetes Recycle Bin
Kubed provides a recycle bin for deleted and/or updated Kubernetes objects. Once activated, any deleted and/or updated object is stored in YAML format in folder mounted inside Kubed pod. This tutorial will show you how to use Kubed to setup a recycle bin for Kubernetes cluster objects.


## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).


## Deploy Kubed
To enable config syncer, you need a cluster config like below.

```yaml
$ cat ./docs/examples/recycle-bin/config.yaml

clusterName: unicorn
notifierSecretName: notifier-config
recycleBin:
  path: /tmp/kubed/trash
  ttl: 168h
  handleUpdates: false
  receivers:
  - notifier: Mailgun
    to:
    - ops@example.com
```

| Key                        | Description                                                                               |
|----------------------------|-------------------------------------------------------------------------------------------|
| `recycleBin.path`          | `Required`. Path to folder where deleted and/or updated objects are stored. |
| `recycleBin.ttl`           | `Required`. Duration for which deleted and/or updated objects are stored before purging. |
| `recycleBin.handleUpdates` | `Optional`. If set to `true`, past version of supported objects are stored when updated. We recommend that you keep this set to `false` on an active cluster. |
| `recycleBin.receiver`      | `Optional`. If set, a notification will be sent when any supported object is deleted and/or updated. To learn how to use various notifiers, please visit [here](/docs/guides/notifiers.md). |
| `clusterName`              | `Optional`. A meaningful identifer for cluster. This cluster name will be prefixed to any notification sent via Email/SMS/Chat so that you can identify the source easily. |

Now, create a Secret with the Kubed cluster config under `config.yaml` key.

```yaml
$ kubectl create secret generic kubed-config -n kube-system \
    --from-file=./docs/examples/recycle-bin/config.yaml
secret "kubed-config" created

# apply app=kubed label to easily cleanup later
$ kubectl label secret kubed-config app=kubed -n kube-system
secret "kubed-config" labeled

$ kubectl get secret kubed-config -n kube-system -o yaml
apiVersion: v1
data:
  config.yaml: bm90aWZpZXJTZWNyZXROYW1lOiBrdWJlZC1ub3RpZmllcgpyZWN5Y2xlQmluOgogIGhhbmRsZV91cGRhdGU6IGZhbHNlCiAgcGF0aDogL3RtcC9rdWJlZAogIHJlY2VpdmVyOgogICAgbm90aWZpZXI6IG1haWxndW4KICAgIHRvOgogICAgLSBvcHNAZXhhbXBsZS5jb20KICB0dGw6IDE2OGgK
kind: Secret
metadata:
  creationTimestamp: 2017-07-26T18:55:54Z
  labels:
    app: kubed
  name: kubed-config
  namespace: kube-system
  resourceVersion: "32920"
  selfLink: /api/v1/namespaces/kube-system/secrets/kubed-config
  uid: 0d3aa21b-7234-11e7-af79-08002738e55e
type: Opaque
```

Now, deploy Kubed operator in your cluster following the steps [here](/docs/setup/install.md). Once the operator pod is running, go to the next section.


## Using Recycle Bin
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

To uninstall Kubed operator, please follow the steps [here](/docs/setup/uninstall.md).

## Next Steps
 - Learn how to use Kubed to take periodic snapshots of a Kubernetes cluster [here](/docs/guides/cluster-snapshot.md).
 - Need to keep some configuration synchronized across namespaces? Try [Kubed config syncer](/docs/guides/config-syncer.md).
 - Want to keep an eye on your cluster with automated notifications? Setup Kubed [event forwarder](/docs/guides/event-forwarder.md).
 - Out of disk space because of too much logs in Elasticsearch or metrics in InfluxDB? Configure [janitors](/docs/guides/janitors.md) to delete old data.
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/ROADMAP.md).
 - Want to hack on Kubed? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
