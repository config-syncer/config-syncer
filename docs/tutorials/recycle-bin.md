> New to Kubed? Please start [here](/docs/tutorials/README.md).

# Kubernetes Recycle Bin
Kubed provides a recycle bin for deleted and/or updated Kubernetes objects. Once activated, any deleted and/or updated object is stored in YAML format in folder mounted inside Kubed pod. This tutorial will show you how to use Kubed to setup a recycle bin for Kubernetes cluster objects.


## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).


## Deploy Kubed
To enable config syncer, you need a cluster config like below.

```yaml
$ cat ./docs/examples/recycle-bin/config.yaml

notifierSecretName: kubed-notifier
recycleBin:
  path: /tmp/kubed
  ttl: 168h
  handle_update: false
  receiver:
    notifier: mailgun
    to:
    - ops@example.com
```

| Key                        | Description                                                                               |
|----------------------------|-------------------------------------------------------------------------------------------|
| `recycleBin.path`          | `Required`. Path to folder where deleted and/or updated objects are stored. |
| `recycleBin.ttl`           | `Required`. Duration for which deleted and/or updated objects are stored before purging. |
| `recycleBin.handle_update` | `Optional`. If set to `true`, past version of supported objects are stored when updated. We recommend that you keep this set to `false` on an active cluster. |
| `recycleBin.receiver`      | `Optional`. If set, a notification will be sent when any supported object is deleted and/or updated. To learn how to use various notifiers, please visit [here](./docs/tutorials/notifiers.md). |

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

Now, deploy Kubed operator in your cluster following the steps [here](/docs/install.md). Once the operator pod is running, go to the next section.


## Using Recycle Bin


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

## Disable Recycle Bin
If you would like to disable this feature, remove the `recyclebin` portion of your Kubed cluster config. Then update the `kubed-config` Secret and restart Kubed operator pod(s).


## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, uninstall Kubed operator following the steps [here](/docs/uninstall.md).

## Next Steps
